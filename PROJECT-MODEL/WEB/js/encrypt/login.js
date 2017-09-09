(function (P) {

    var global = {
        url: '/gotrix/',
        nonSessionPage: 'index.html',
        noSession: function (reqParams) {
            P.alert('用户会话已过期，请重新登陆。');
        }
    };

    P.global = global;

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

    P.rand = function () {
        var n = new BigInteger(N);
        var n1 = n.subtract(BigInteger.ONE);
        var r = new BigInteger(n.bitLength(), get_rng());
        return r.mod(n1).add(BigInteger.ONE);
    }

    P.secretKey = function (a) {
        var before = new Date();

        var curve = get_curve();
        var G = get_G(curve);
        var S = G.multiply(a);

        var after = new Date();
        console.info('耗时：' + (after - before) + 'ms');
        return S;
    }

    P.publicKey = function (a, x, y) {
        var before = new Date();

        var curve = get_curve();
        var G = new ECPointFp(curve, curve.fromBigInteger(new BigInteger(x, 16)), curve.fromBigInteger(new BigInteger(y, 16)));
        var P = G.multiply(a);

        var after = new Date();
        console.info('耗时：' + (after - before) + 'ms');
        return P;
    }

    P.successFunc = function (d) {
        P.alert(d.msg);
    }

    P.errorFunc = function (d) {
        P.alert(d.msg || '网络异常');
    }

    function writeSession(session, password) {
        P.ls('session', {
            session: session,
            password: password
        });
    }

    function getSession() {
        return P.ls('session');
    }

    P.extend = function () {
        /*
         　　*target被扩展的对象
         　　*length参数的数量
         　　*deep是否深度操作
         　　*/
        var options, name, src, copy, copyIsArray, clone,
            target = arguments[0] || {},
            i = 1,
            length = arguments.length,
            deep = false;

        // target为第一个参数，如果第一个参数是Boolean类型的值，则把target赋值给deep
        // deep表示是否进行深层面的复制，当为true时，进行深度复制，否则只进行第一层扩展
        // 然后把第二个参数赋值给target
        if (typeof target === "boolean") {
            deep = target;
            target = arguments[1] || {};

            // 将i赋值为2，跳过前两个参数
            i = 2;
        }

        // target既不是对象也不是函数则把target设置为空对象。
        if (typeof target !== "object" && !jQuery.isFunction(target)) {
            target = {};
        }

        // 如果只有一个参数，则把jQuery对象赋值给target，即扩展到jQuery对象上
        if (length === i) {
            target = this;

            // i减1，指向被扩展对象
            --i;
        }

        // 开始遍历需要被扩展到target上的参数

        for (; i < length; i++) {
            // 处理第i个被扩展的对象，即除去deep和target之外的对象
            if ((options = arguments[i]) != null) {
                // 遍历第i个对象的所有可遍历的属性
                for (name in options) {
                    // 根据被扩展对象的键获得目标对象相应值，并赋值给src
                    src = target[name];
                    // 得到被扩展对象的值
                    copy = options[name];

                    // 这里为什么是比较target和copy？不应该是比较src和copy吗？
                    if (target === copy) {
                        continue;
                    }

                    // 当用户想要深度操作时，递归合并
                    // copy是纯对象或者是数组
                    if (deep && copy && ( jQuery.isPlainObject(copy) || (copyIsArray = jQuery.isArray(copy)) )) {
                        // 如果是数组
                        if (copyIsArray) {
                            // 将copyIsArray重新设置为false，为下次遍历做准备。
                            copyIsArray = false;
                            // 判断被扩展的对象中src是不是数组
                            clone = src && jQuery.isArray(src) ? src : [];
                        } else {
                            // 判断被扩展的对象中src是不是纯对象
                            clone = src && jQuery.isPlainObject(src) ? src : {};
                        }

                        // 递归调用extend方法，继续进行深度遍历
                        target[name] = jQuery.extend(deep, clone, copy);

                        // 如果不需要深度复制，则直接把copy（第i个被扩展对象中被遍历的那个键的值）
                    } else if (copy !== undefined) {
                        target[name] = copy;
                    }
                }
            }
        }

        // 原对象被改变，因此如果不想改变原对象，target可传入{}
        return target;
    };

    P.ajax = function (ajaxParam) {
        /*
         * 创建XMLHttpRequest对象
         */
        var xmlHttp;
        try {
            if (window.plus && window.plus.net && window.plus.net.XMLHttpRequest) {
                // HBuilder H5+ 环境下的跨域请求
                xmlHttp = new plus.net.XMLHttpRequest();
            } else {
                // Firefox, Opera 8.0+, Safari
                xmlHttp = new XMLHttpRequest();
            }
        }
        catch (e) {
            try {// Internet Explorer
                xmlHttp = new ActiveXObject("Msxml2.XMLHTTP");
            }
            catch (e) {
                try {
                    xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
                }
                catch (e) {
                }
            }
        }

        /*
         * 服务器向浏览器响应请求
         *
         * readyState 属性表示Ajax请求的当前状态。它的值用数字代表。
         * 0 代表未初始化。 还没有调用 open 方法
         * 1 代表正在加载。 open 方法已被调用，但 send 方法还没有被调用
         * 2 代表已加载完毕。send 已被调用。请求已经开始
         * 3 代表交互中。服务器正在发送响应
         * 4 代表完成。响应发送完毕

         * 常用状态码及其含义：
         * 404 没找到页面(not found)
         * 403 禁止访问(forbidden)
         * 500 内部服务器出错(internal service error)
         * 200 一切正常(ok)
         * 304 没有被修改(not modified)(服务器返回304状态，表示源文件没有被修改 )
         */
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4) {
                if (xmlHttp.status == 200 || xmlHttp.status == 304) {
                    ajaxParam.success && ajaxParam.success(xmlHttp.responseText);
                } else {
                    ajaxParam.error && ajaxParam.error(xmlHttp.responseText);
                }
            }
        }

        /*
         * 3    浏览器与服务器建立连接
         *
         * xhr.open(method, url, asynch);
         *         * 与服务器建立连接使用
         *         * method：请求类型，类似 “GET”或”POST”的字符串。
         *         * url：路径字符串，指向你所请求的服务器上的那个文件。请求路径
         *         * asynch：表示请求是否要异步传输，默认值为true(异步)。
         */
        xmlHttp.open(ajaxParam.type, ajaxParam.url, ajaxParam.async);

        //如果是POST请求方式，设置请求首部信息
        xmlHttp.setRequestHeader("Content-type", ajaxParam.contentType);


        /*
         * 4    浏览器向服务器发送请求
         *
         *     send()方法：
         *         * 如果浏览器请求的类型为GET类型时，通过send()方法发送请求数据，服务器接收不到
         *         * 如果浏览器请求的类型为POST类型时，通过send()方法发送请求数据，服务器可以接收
         */
        // var param = '';
        // for (var key in data) {
        //     param += key + '=' + (key == 'data' ? JSON.stringify(data[key]) : data[key]) + "&";
        // }
        xmlHttp.send(ajaxParam.data);        //xhr.send(null);
    }

    // 请求队列
    P.reqQueue = [];

    P.reqing = false;

    P.req = function (data, success, error, config) {

        // 插入队列之首
        P.reqQueue.unshift({
            data: data,
            success: success,
            error: error,
            config: config
        });

        P._req();

    }

    P._req = function () {

        // 有请求正在执行，则直接返回。
        if (P.reqing) {
            return;
        }

        // 取队列最末数据
        var reqParams = P.reqQueue.pop();
        if (!reqParams) {
            return;
        }

        // 设置请求正在执行
        P.reqing = true;

        data = reqParams.data;
        success = reqParams.success;
        error = reqParams.error;
        config = reqParams.config;

        // 在某些情况下需要重新调用该接口
        var recallReq = function () {
            P.req(data, success, error, config);
        };

        // 初始化默认参数：
        var successFunc = success || P.successFunc;
        var errorFunc = error || P.errorFunc;
        var ajaxData = P.extend({}, data);
        var ajaxConfig = P.extend({}, config);

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
            token = ajaxData.TOKEN;
            content = Aes.Ctr.encrypt(JSON.stringify(ajaxData), ajaxConfig.aesPass, 256);
        } else if ((session = getSession()) && session.session && session.password) {
            token = session.session;
            content = Aes.Ctr.encrypt(JSON.stringify(ajaxData), session.password, 256);
        } else {
            P.reqing = false;
            global.noSession(reqParams);
            return;
        }

        if (ajaxConfig.blank) {
            // TODO 测试一下这个功能
            var form = document.createElement('form');
            form.style.display = 'none';
            form.action = global.url;
            form.method = 'POST';
            form.enctype = 'text/plain';
            form.innerHTML = '<input name="' + token + '" value="' + content + '" />';
            form.commit();
            return;
        } else {
            ajaxData = token + '=' + content;
        }

        // 组装请求参数：
        var ajaxParam = P.extend({
            async: true,
            url: global.url,
            type: 'POST',
            dataType: 'text',
            contentType: 'text/plain',
            data: ajaxData,
            jsonp: false,
            success: function (text) {
                // 返回成功，继续执行其它请求
                P.reqing = false;
                P._req();

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
                } else if (d.status === 1000 && d.msg && d.msg.indexOf('用户会话已过期') === 0) { // 用户 Session 已过期
                    console.error(d);
                    global.noSession(reqParams);
                } else {
                    console.error(d);
                    errorFunc(d);
                }

                // 如果有 intervalTime 这个参数，则相隔 intervalTime 之后持续向服务器请求：
                if (ajaxConfig.intervalTime) {
                    setTimeout(recallReq, ajaxConfig.intervalTime);
                }
            },
            error: function (text) {
                // 发生错误，继续执行其它请求。
                P.reqing = false;
                P._req();

                console.error(text);
                errorFunc(text);
            }
        }, ajaxConfig);

        // 向服务器发送请求：
        if (ajaxConfig.ajaxFunc) {
            ajaxConfig.ajaxFunc(ajaxParam);
        } else {
            P.ajax(ajaxParam);
        }
    }

    P.codeLogin = function (code, parent, success, error) {
        var rand = P.rand();
        var S = P.secretKey(rand);

        $.self({
            func: 9,
            code: code,
            parent: parent,
            x: S.getX().toBigInteger().toString(16),
            y: S.getY().toBigInteger().toString(16)
        }, function (data) {
            var session = data.data.session;
            var PK = P.publicKey(rand, data.data.x, data.data.y);
            var password = PK.getX().toBigInteger().add(PK.getY().toBigInteger()).toString(16);
            writeSession(session, password);
            // success && success(data.data);
            window.location = (window.location + '').replace(/com\/?\d*/, 'com/' + data.data.id + '/').replace(/code=.*&/, '');
        }, error);
    }

    P.UUIDLogin = function (uuid, success, error) {
        var rand = P.rand();
        var S = P.secretKey(rand);

        P.self({
            func: 14,
            uuid: uuid,
            x: S.getX().toBigInteger().toString(16),
            y: S.getY().toBigInteger().toString(16)
        }, function (data) {
            var session = data.data.session;
            var PK = P.publicKey(rand, data.data.x, data.data.y);
            var password = PK.getX().toBigInteger().add(PK.getY().toBigInteger()).toString(16);
            writeSession(session, password);
            success && success();
        }, error)

    };

    P.loginIn = function (username, password, success, error) {
        var tokenSalt = "!QAZXSW@#EDCVFR$"
        var token = hex_sha512(username + tokenSalt + password).substring(32, 96);

        var bcryptSalt = "$2a$08$CTxJ2aIST2x1EAcLg934ie";
        var aesPassword = dcodeIO.bcrypt.hashSync(password, bcryptSalt).replace(bcryptSalt, "");

        P.req({
            func: 'GetSalt',
            TOKEN: token
        }, getSession, error, {
            aesPass: aesPassword
        });

        function getSession(data) {
            var salt = data.data;
            var aesPassword = dcodeIO.bcrypt.hashSync(password + salt, salt).replace(salt, "");

            var rand = P.rand();
            var S = P.secretKey(rand);
            P.req({
                func: 'LoginIn',
                TOKEN: token,
                x: S.getX().toBigInteger().toString(16),
                y: S.getY().toBigInteger().toString(16)
            }, function (data) {
                var session = data.data.session;
                var PK = P.publicKey(rand, data.data.x, data.data.y);
                var password = PK.getX().toBigInteger().add(PK.getY().toBigInteger()).toString(16);
                writeSession(session, password);
                success && success();
            }, error, {
                aesPass: aesPassword
            });
        }

    }

    P.register = function (username, password, success, error) {
        var tokenSalt = "!QAZXSW@#EDCVFR$"
        var token = hex_sha512(username + tokenSalt + password).substring(32, 64);

        var bcryptSalt = "$2a$08$CTxJ2aIST2x1EAcLg934ie";
        var pass = dcodeIO.bcrypt.hashSync(password, bcryptSalt).replace(bcryptSalt, "");

        var salt = dcodeIO.bcrypt.genSaltSync(10);
        var key = dcodeIO.bcrypt.hashSync(password + salt, salt).replace(salt, "");

        var aesPassword = token.substring(0, 31);
        var rand = P.rand();
        var S = P.secretKey(rand);
        P.req({
            func: 13,
            TOKEN: token,
            pass: pass,
            salt: salt,
            key: key
        }, success, error, {
            aesPass: aesPassword
        })
    }

    P.self = function (data, success, error) {
        var aesPassword = hex_sha512(dcodeIO.bcrypt.genSaltSync(1)).substring(32, 64);
        P.req(P.extend(data, {
            TOKEN: aesPassword
        }), success, error, {
            aesPass: aesPassword.substring(0, 31)
        })
    }

})(window.P || window.jQuery || window.mui);