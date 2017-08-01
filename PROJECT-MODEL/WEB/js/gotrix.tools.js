;
(function (w) {

    var P = w.P || (w.P = {});

    /**
     * 扩展 String 的 from 方法。
     */
    String.prototype.from = function (obj, filter) {
        var g_reg = /\*\{([\w\.])*?(#[^\}]*)?\}/g;
        var reg1 = /\*\{(\w*?)(#([^\}]*))?\}/;
        var reg2 = /\*\{([\w\.]*?)(#([^\}]*))?\}/;
        obj = obj || [];
        if (obj.constructor == Array) {
            var strs = [], d, properties, resultStr, exs, replace;
            for (var i = 0; i < obj.length; i++) {
                d = obj[i] || {};
                if (filter && filter(i, d) === false) {
                    continue;
                }
                strs.push(this.replace(g_reg, function (s) {
                    exs = reg1.exec(s);
                    if (exs) {
                        replace = d[exs[1]];
                        return (replace === undefined || replace === null) ? (exs[3] || '') : replace;
                    }
                    exs = reg2.exec(s);
                    properties = exs[1].split('.');
                    resultStr = d;
                    for (var i = 0; i < properties.length; i++) {
                        resultStr = resultStr[properties[i]];
                        if (!resultStr) {
                            return exs[3] || '';
                        }
                    }
                    return resultStr;
                }));
            }
            return strs.join('');
        } else if (obj.constructor == Object) {
            return this.from([obj], filter);
        }
    }

    /**
     * 周期执行事件
     */
    w.requestAnimFrame = window.requestAnimationFrame || window.mozRequestAnimationFrame || window.webkitRequestAnimationFrame || window.msRequestAnimationFrame || window.oRequestAnimationFrame || function (callback) {
            return setTimeout(callback, 1000 / 60);
        };

    /**
     * 周期执行事件
     */
    w.cancelAnimFrame = window.cancelAnimationFrame || window.webkitCancelAnimationFrame || window.webkitCancelRequestAnimationFrame || window.mozCancelRequestAnimationFrame || window.oCancelRequestAnimationFrame || window.msCancelRequestAnimationFrame || clearTimeout;

    /**
     * 从 url 中获取出所有的参数
     */
    P.getParam = function (name) {
        var param = {};
        var reg = /(?:\?|&)(\w*?)=([^&\=]*?)(?=&|$)/g;
        var href = window.location.href;
        var strs = reg.exec(href);
        while (strs) {
            param[strs[1]] = strs[2];
            strs = reg.exec(href);
        }
        return name ? param[name] : param;
    }

    /**
     * 判断一个元素是否在这个数组中
     */
    P.inArray = function (d, arr) {
        if (typeof arr == 'object' && arr.constructor == Array) {
            for (var i = 0, len = arr.length; i < len; i++) {
                if (arr[i] == d) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * 获取当前时间，或某个时间的详细信息。
     */
    P.getDate = function (date) {
        date = date || new Date();
        return {
            year: date.getFullYear(),
            month: date.getMonth() + 1,
            date: date.getDate(),
            day: date.getDay(),
            dayOfWeek: ['日', '一', '二', '三', '四', '五', '六'][date.getDay()],
            hour: date.getHours(),
            minute: date.getMinutes(),
            second: date.getSeconds(),
            millisecond: date.getMilliseconds()
        }
    }

    /**
     * 获取当前日期
     */
    P.getCurDate = function (date, join) {
        var d = P.getDate(date);
        join = join || '';
        return d.year + join + d.month + join + d.date;
    }

    /**
     * 图片修复
     */
    w.imgSizer = {
        Config: {
            imgCache: [],
            spacer: "/path/to/your/spacer.gif"
        },
        collate: function (aScope) {
            var isOldIE = (document.all && !window.opera && !window.XDomainRequest) ? 1 : 0;
            if (isOldIE && document.getElementsByTagName) {
                var c = imgSizer;
                var imgCache = c.Config.imgCache;

                var images = (aScope && aScope.length) ? aScope : document.getElementsByTagName("img");
                for (var i = 0; i < images.length; i++) {
                    images[i].origWidth = images[i].offsetWidth;
                    images[i].origHeight = images[i].offsetHeight;

                    imgCache.push(images[i]);
                    c.ieAlpha(images[i]);
                    images[i].style.width = "100%";
                }

                if (imgCache.length) {
                    c.resize(function () {
                        for (var i = 0; i < imgCache.length; i++) {
                            var ratio = (imgCache[i].offsetWidth / imgCache[i].origWidth);
                            imgCache[i].style.height = (imgCache[i].origHeight * ratio) + "px";
                        }
                    });
                }
            }
        },
        ieAlpha: function (img) {
            var c = imgSizer;
            if (img.oldSrc) {
                img.src = img.oldSrc;
            }
            var src = img.src;
            img.style.width = img.offsetWidth + "px";
            img.style.height = img.offsetHeight + "px";
            img.style.filter = "progid:DXImageTransform.Microsoft.AlphaImageLoader(src='" + src + "', sizingMethod='scale')"
            img.oldSrc = src;
            img.src = c.Config.spacer;
        },
        resize: function (func) {
            var oldonresize = window.onresize;
            if (typeof window.onresize != 'function') {
                window.onresize = func;
            } else {
                window.onresize = function () {
                    if (oldonresize) {
                        oldonresize();
                    }
                    func();
                }
            }
        }
    }

})(window);