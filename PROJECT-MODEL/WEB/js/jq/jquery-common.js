(function () {

    var global = {
        url: '/gotrix/',
        nonSessionPage: 'index.html'
    }

    var Q = '1461501637330902918203684832716283019653785059327';
    var A = '1461501637330902918203684832716283019653785059324';
    var B = '163235791306168110546604919403271579530548345413';
    var GX = '425826231723888350446541592701409065913635568770';
    var GY = '203520114162904107873991457957346892027982641970';
    var N = '1461501637330902918203687197606826779884643492439';

    var rng;

    function get_rng() {
        return rng || (rng = new SecureRandom());
    }

    function get_curve() {
        return new ECCurveFp(new BigInteger(Q), new BigInteger(A), new BigInteger(B));
    }

    function get_G(curve) {
        return new ECPointFp(curve, curve.fromBigInteger(new BigInteger(GX)), curve.fromBigInteger(new BigInteger(GY)));
    }

    $.rand = function () {
        var n = new BigInteger(N);
        var n1 = n.subtract(BigInteger.ONE);
        var r = new BigInteger(n.bitLength(), get_rng());
        return r.mod(n1).add(BigInteger.ONE);
    }

    $.secretKey = function (a) {
        var before = new Date();

        var curve = get_curve();
        var G = get_G(curve);
        var S = G.multiply(a);

        var after = new Date();
        console.info('耗时：' + (after - before) + 'ms');
        return S;
    }

    $.publicKey = function (a, x, y) {
        var before = new Date();

        var curve = get_curve();
        var G = new ECPointFp(curve, curve.fromBigInteger(new BigInteger(x, 16)), curve.fromBigInteger(new BigInteger(y, 16)));
        var P = G.multiply(a);

        var after = new Date();
        console.info('耗时：' + (after - before) + 'ms');
        return P;
    }

    $.successFunc = function (d) {
        P.alert(d.msg);
    }

    $.errorFunc = function (d) {
        P.alert(d.msg || '网络异常');
    }

    function writeSession(session, password) {
        P.ss('session', {
            session: session,
            password: password
        });
    }

    function getSession() {
        return P.ss('session');
    }

    $.req = function (data, success, error, config) {
        // 初始化默认参数：
        var successFunc = success || $.successFunc;
        var errorFunc = error || $.errorFunc;
        var ajaxData = $.extend({}, data);
        var ajaxConfig = $.extend({}, config);

        // 对 ajaxData 进行去除特殊字符处理：
        for (var key in ajaxData) {
            if (typeof ajaxData[key] == 'string' && key != 'password') {
                ajaxData[key] = ajaxData[key].replace(/\\/g, '\\\\').replace(/"/g, '\\\"').replace(/\r\n/g, '\\n').replace(/\n/g, '\\n');
            }
        }

        // 对数据进行加密处理：
        var session;
        var token;
        var content;
        if (ajaxConfig.aesPass) {
            token = ajaxData.token;
            content = Aes.Ctr.encrypt(JSON.stringify(ajaxData), ajaxConfig.aesPass, 256);
        } else if ((session = getSession()) && session.session && session.password) {
            token = session.session;
            content = Aes.Ctr.encrypt(JSON.stringify(ajaxData), session.password, 256);
        } else {
            P.alert('您尚未登录');
            return;
        }

        if (ajaxConfig.blank) {
            $('<form style="display:none;" action="' + global.url + '" method="POST" enctype="text/plain"><input name="' + token + '" value="' + content + '" /></form>').appendTo('body').submit();
            return;
        } else {
            ajaxData = token + '=' + content;
        }

        // 组装请求参数：
        var ajaxParam = $.extend({
            async: true,
            url: global.url,
            type: 'POST',
            dataType: 'text',
            contentType: 'text/plain',
            data: ajaxData,
            jsonp: false,
            success: function (text) {
                // 解密服务器返回的内容：
                var d;
                try {
                    d = JSON.parse(Aes.Ctr.decrypt(text, ajaxConfig.aesPass || getSession().password, 256));
                } catch (e) {
                    try {
                        d = JSON.parse(text.replace(/\n/, '\\\n'));
                    } catch (e1) {
                        errorFunc(text);
                        return;
                    }
                }

                // 解密成功之后，判断服务器返回的状态，状态为 0 调用成功的回调函数，状态不为 0 调用失败的回调函数：
                if (d.status === 0) {
                    console.info(d);
                    successFunc(d);
                } else {
                    console.error(d);
                    errorFunc(d);
                }

                // 如果有 intervalTime 这个参数，则相隔 intervalTime 之后持续向服务器请求：
                if (ajaxConfig.intervalTime) {
                    setTimeout(function () {
                        $.req(data, success, error, config);
                    }, ajaxConfig.intervalTime);
                }
            },
            error: function (text) {
                console.error(text);
                errorFunc(text);
            }
        }, ajaxConfig);

        // 向服务器发送请求：
        $.ajax(ajaxParam);
    }

    $.codeLogin = function (code, parent, success, error) {
        var rand = $.rand();
        var S = $.secretKey(rand);

        $.self({
            func: 9,
            code: code,
            parent: parent,
            x: S.getX().toBigInteger().toString(16),
            y: S.getY().toBigInteger().toString(16)
        }, function (data) {
            var session = data.data.session;
            var PP = $.publicKey(rand, data.data.x, data.data.y);
            var password = PP.getX().toBigInteger().add(PP.getY().toBigInteger()).toString(16);
            writeSession(session, password);
            // success && success(data.data);
            window.location = (window.location + '').replace(/com\/?\d*/, 'com/' + data.data.id + '/').replace(/code=.*&/, '');
        }, error);
    }

    $.loginIn = function (username, password, success, error) {
        var tokenSalt = "!QAZXSW@#EDCVFR$"
        var token = hex_sha512(username + tokenSalt + password).substring(32, 96);

        var bcryptSalt = "$2a$08$CTxJ2aIST2x1EAcLg934ie";
        var aesPassword = dcodeIO.bcrypt.hashSync(password, bcryptSalt).replace(bcryptSalt, "");

        $.req({
            func: 10,
            token: token
        }, getSession, error, {
            aesPass: aesPassword
        });

        function getSession(data) {
            var salt = data.data;
            var aesPassword = dcodeIO.bcrypt.hashSync(password + salt, salt).replace(salt, "");

            var rand = $.rand();
            var S = $.secretKey(rand);
            $.req({
                func: 11,
                token: token,
                x: S.getX().toBigInteger().toString(16),
                y: S.getY().toBigInteger().toString(16)
            }, function (data) {
                var session = data.data.session;
                var P = $.publicKey(rand, data.data.x, data.data.y);
                var password = P.getX().toBigInteger().add(P.getY().toBigInteger()).toString(16);
                writeSession(session, password);
                success();
            }, error, {
                aesPass: aesPassword
            });
        }

    }

    $.register = function (username, password, success, error) {
        var tokenSalt = "!QAZXSW@#EDCVFR$"
        var token = hex_sha512(username + tokenSalt + password).substring(32, 64);

        var bcryptSalt = "$2a$08$CTxJ2aIST2x1EAcLg934ie";
        var pass = dcodeIO.bcrypt.hashSync(password, bcryptSalt).replace(bcryptSalt, "");

        var salt = dcodeIO.bcrypt.genSaltSync(10);
        var key = dcodeIO.bcrypt.hashSync(password + salt, salt).replace(salt, "");

        var aesPassword = token.substring(0, 31);
        var rand = $.rand();
        var S = $.secretKey(rand);
        $.req({
            func: 13,
            token: token,
            pass: pass,
            salt: salt,
            key: key
        }, success, error, {
            aesPass: aesPassword
        })
    }

    $.self = function (data, success, error) {
        var aesPassword = hex_sha512(dcodeIO.bcrypt.genSaltSync(1)).substring(32, 64);
        $.req($.extend(data, {
            token: aesPassword
        }), success, error, {
            aesPass: aesPassword.substring(0, 31)
        })
    }

    $.addFloat = function (i) {
        var strs = ['<div class="float">',

            '<a href="index.html"><img src="img/index_float1.png" /><img class="hide" src="img/index_float1_active.png" /><p>首页</p></a>',

            '<a href="shopping.html"><img src="img/index_float2.png" /><img class="hide" src="img/index_float2_active.png" /><p>购物车</p></a>',

            '<a href="store.html"><img src="img/index_float3.png" /><img class="hide" src="img/index_float3_active.png" /><p>我的人脉</p></a>',

            '<a href="person.html"><img src="img/index_float4.png" /><img class="hide" src="img/index_float4_active.png" /><p>会员</p></a>',

            '</div>'];

        $('body').append(strs.join('')).find('.float a:eq(' + (i || 0) + ')').find('img:first').hide().next().show();
    }

    $.waterfall = function (data, success, error) {

        var curPage = 0, pageNumber = 5;
        var loading = false;
        var allLoaded = false;

        function loadData() {
            if (allLoaded || loading) {
                return;
            }
            loading = true;
            $.self($.extend(data, {
                start: pageNumber * curPage,
                end: pageNumber
            }), function (data) {
                if (data && data.data && data.data.list) {
                    if (!data.data.list.length) {
                        allLoaded = true;
                    }
                    success && success(data);
                }
                loading = false;
            }, function (data) {
                error && error(data);
                loading = false;
            });
            curPage++;
        }

        loadData();

        $(window).scroll(function (e) {
            if ($(window).scrollTop() >= $(document).height() - $(window).height()) {
                loadData();
            }
        });
    }

    $(function () {

        if (!window.rem) {
            rem = Math.min(600, $(window).width()) / 10;
            $('html,body').css({
                'font-size': rem + 'px',
                'max-width': '600px',
                'margin': '0 auto'
            });
        }

        window.ready && window.ready();

        // $.RP = P.getParam();
        // $.userid = /\/(\d+)\//.exec(window.location);
        // $.userid = $.userid ? $.userid[1] : 0;
        //
        // if (window.beforeLogin) {
        //     window.beforeLogin();
        // } else if (!getSession()) {
        //     if ($.RP.code) {
        //         $.codeLogin($.RP.code, $.userid, window.ready);
        //     } else {
        //         window.location = 'https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx32e598477c7d1ef8&redirect_uri=' + window.location + '&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect';
        //     }
        // } else {
        //     window.ready && window.ready();
        // }

    });

})()
