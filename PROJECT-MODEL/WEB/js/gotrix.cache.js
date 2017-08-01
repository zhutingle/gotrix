;
(function (w) {

    /**
     * IE7及以下没有JSON这个内置对象，所以给它创建一个
     */
    if (!w.JSON) {
        JSON = {
            parse: function (str) {
                return eval('(' + str + ')');
            },
            stringify: function (value) {
                if (!value && value !== false && value != 0) {
                    return "null";
                }
                var a = [];
                var type = typeof value;
                if (type == "number") {
                    a.push(value);
                } else if (type == 'string') {
                    a.push('"' + value.replace(/\\/g, '\\\\').replace(/"/g, '\\\"').replace(/\r\n/g, '\\n').replace(/\n/g, '\\n') + '"');
                } else if (type != 'object') {
                    a.push(value);
                } else {
                    if (value.constructor == Date) {
                        a.push(value);
                    } else if (value.constructor == Array) {
                        a.push("[");
                        var i = 0;
                        for (; i < value.length; i++) {
                            a.push(JSON.stringify(value[i]));
                            a.push(",");
                        }
                        if (i > 0) {
                            a.pop();
                        }
                        a.push("]");
                    } else {
                        a.push("{");
                        var key = null;
                        for (key in value) {
                            a.push('"' + key + '":');
                            a.push(JSON.stringify(value[key]));
                            a.push(",");
                        }
                        if (key) {
                            a.pop();
                        }
                        a.push("}");
                    }
                }
                return a.join("");
            }
        };
    }

    /**
     * 封装 LocalStorage 和 SessionStorage
     */
    w.P = (function () {
        var storages = {
            ls: w.localStorage,
            ss: w.sessionStorage
        }
        for (var key in storages) {
            storages[key] = (function (key, storage) {
                var supported = typeof storage != 'undefined' && typeof storage != 'unknown';

                function func(tempKey, value) {
                    if (!supported) {
                        return false;
                    }
                    if (tempKey === undefined) {
                        if (value && value.test) {
                            for (var k in storage) {
                                value.test(k) && storage.removeItem(k);
                            }
                        } else {
                            for (var k in storage) {
                                storage.removeItem(k);
                            }
                        }
                        return true;
                    } else if (value === undefined) {
                        try {
                            return JSON.parse(storage.getItem(tempKey));
                        } catch (e) {
                            return storage.getItem(tempKey);
                        }
                    } else if (value == null) {
                        return storage.removeItem(tempKey);
                    } else {
                        return storage.setItem(tempKey, typeof value == 'object' ? JSON.stringify(value) : value);
                    }
                }

                func.supported = supported;
                return func;
            })(key, storages[key]);
        }

        return storages;

    })();

    /**
     * 封装弹出提示框
     *
     * @param msg
     */
    P.alert = function (msg) {
        var alert_content = document.getElementById('alert_content');
        if (!alert_content) {
            alert_content = document.createElement('div');
            alert_content.id = 'alert_content';
            alert_content.setAttribute('style', 'position: fixed;width: 100%;height: 100%;top: 0;left: 0;text-align: center;z-index: 9999;');
            alert_content.innerHTML = '<div style="position: absolute;width: 6rem;left: 50%;bottom: 3rem;margin: 0 0 0 -3.5rem;display: block;line-height: 0.6rem;font-size: 0.5rem;background: rgba(219, 87, 51, 1);padding: 0.5rem;font-weight: 700;border-radius: 0.5rem;color: white;box-shadow: 0px 9px 0px rgba(219, 31, 5, 1), 0px 9px 25px rgba(0, 0, 0, .7);"></div>';
            document.body.appendChild(alert_content);
        }
        alert_content.style.display = 'block';
        alert_content.childNodes[0].innerHTML = msg;

        setTimeout(function () {
            alert_content.style.display = 'none';
        }, 2000);
    }

    /**
     * 输出日志，级别为 info
     */
    P.consoleInfo = function (msg) {
        if (window.GOTRIX_INFO && window.console && console.info) {
            window.console.info(msg);
        }
    }

    /**
     * 输出日志，级别为 error
     */
    P.consoleError = function (msg) {
        if (window.GOTRIX_ERROR && window.console && console.error) {
            window.console.error(msg);
        }
    }

    /**
     * 封装 Cache
     */
    P.c = (function () {

        var F_EXISTS = typeof F != 'undefined';

        setTimeout(function () {
            loadUrl('js/gotrix.cache.js', undefined, 3, true);
        }, 2000);

        function loadUrl(url, callback, time, noExec) {

            // callback 为 true 时，自动以“,”号拆分 url，并依次导入。
            if (callback === true) {
                for (var urls = url.split(','), i = 0; i < urls.length; i++) {
                    if (urls[i] && loadUrl(urls[i], undefined, 3) === false) {
                        P.alert('加载失败，请刷新后重试。<a style="color:blue;text-decoration:underline;" href="javascript:window.location.reload();" />刷新</a>');
                        break;
                    }
                }
                return;
            }

            var content;
            if (url && F_EXISTS && F[url] == P.ls(url + '_v') && (content = P.ls(url))) {
                if (evalFile(url, content, noExec) === false) { // 若缓存的文件有问题，则再次从网络中获取
                    P.consoleInfo("From cache error:" + url + "    timeStamp:" + F[url]);
                } else {
                    P.consoleInfo("From cache:" + url + "   timeStamp:" + F[url]);
                    return content;
                }
            }
            if (url) {
                P.consoleInfo("Load from Internet:" + url + (F_EXISTS ? '    timeStamp:' + F[url] + '    localStamp:' + P.ls(url + '_v') : ''));
                var o = /MSIE [678]/.test(navigator.userAgent) ? new ActiveXObject('Microsoft.XMLHTTP') : new XMLHttpRequest();
                o.open('get', url + '?timeStamp=' + w.VERSION, false);
                o.send(null);
                if (evalFile(url, o.responseText, noExec) === false) {
                    if (time > 1) {
                        return loadUrl(url, callback, time - 1);
                    } else {
                        return false;
                    }
                } else {
                    P.ls(url, o.responseText);
                    F_EXISTS && F[url] && P.ls(url + '_v', F[url]);
                    return o.responseText;
                }
            }
        }

        function evalFile(file, string, noExec) {
            if (noExec) {
                return;
            }
            try {
                if (/\.js$/.test(file)) {
                    w.eval.call(w, string);
                    // eval(string);
                } else if (/\.css$/.test(file)) {
                    if (/^</.test(string)) {
                        return false;
                    } else {
                        document.write('<style type="text/css">' + string + '</style>');
                    }
                }
            } catch (e) {
                console.error(e);
                return false;
            }
        }

        return loadUrl;

    })();

    /**
     * 封装图片自动加载函数
     *
     * @param dom
     * @param img
     */
    P.img = function (dom) {
        var execs;
        if (dom && /^data:image\/gif;/.test(dom.src) && (execs = /cache="([^"]*?)"/.exec(dom.outerHTML)) && execs[1]) {
            if (/http/.test(execs[1])) {
                dom.src = execs[1];
            } else {
                dom.src = 'data:image/' + execs[2] + ';base64,' + P.c(execs[1] + '.cache');
            }
        }
    }

    P.c('css/reset.css,js/gotrix.tools.js,js/project.js,' + (window.RES || ''), true);

    return true;

})(window);
