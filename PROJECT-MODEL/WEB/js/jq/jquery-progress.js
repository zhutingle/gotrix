;
(function() {

	var plugin = 'progress';

	$.fn[plugin] = function(param, ext1, ext2) {
		if (typeof param == 'string') {
			var returnValue;
			this.each(function() {
				var data = $(this).data(plugin);
				if (data && data[param] && (returnValue = data[param](ext1, ext2)) !== undefined) {
					return false;
				}
			});
			return returnValue === undefined ? this : returnValue;
		} else {
			return this.each(function() {
				$(this).data(plugin, init.init(this, $.extend({}, defaultParam, param)));
				render.render(this);
			});
		}
	}

	var defaultParam = {
		type : 'circle',
		value : 25,
		total : 100,

		background : '#12415F',
		frontColor : '#FFF61D',
		darkColor : '#1F4E6C',

		onSuccess : function(param) {

		}

	}

	var publicFunc = {}

	var init = {
		init : function(dom, param) {
			param.$this = $(dom);
			param.w = $(dom).width();
			param.h = $(dom).height();
			return param;
		}
	}

	var render = {
		render : function(dom) {
			var param = $(dom).data('progress');
			render[param.type] && render[param.type](param);
		},
		circle : function(param) {
			var $this = param.$this, w = param.w, h = param.h, background = param.background, frontColor = param.frontColor, darkColor = param.darkColor;
			var l = Math.min(w, h);
			$this.css('position', 'relative');
			var markStr = '<div style="position:absolute;width:' + (l - 8) + 'px;height:' + (l - 8) + 'px;margin:4px 0 0 4px;background:' + background + ';border-radius:2000px;"></div>';
			var strs = [];
			strs.push('<div class="back" style="position:absolute;width:' + l + 'px;height:' + l + 'px;background:' + darkColor + ';box-shadow:0px 0px 0px 1px ' + darkColor + ';border-radius:2000px;">' + markStr + '</div>');
			strs.push('<div class="right" style="position:absolute;width:' + l + 'px;height:' + l + 'px;background:' + frontColor + ';box-shadow:0px 0px 0px 1px ' + darkColor + ';border-radius:2000px;clip:rect(0,' + l + 'px,' + l + 'px,' + (l / 2) + 'px);">' + markStr + '</div>');
			strs.push('<div class="left" style="position:absolute;width:' + l + 'px;height:' + l + 'px;background:' + frontColor + ';box-shadow:0px 0px 0px 1px ' + darkColor + ';border-radius:2000px;clip:rect(0,' + (l / 2) + 'px,' + l + 'px,0px);">' + markStr + '</div>');
			$this.html(strs.join(''));
			render.circleSet(param);
		},
		circleSet : function(param) {
			var ratio = param.value / param.total, $this = param.$this;
			if (ratio > 0.5) {
				$this.find('.left').css('background-color', param.frontColor).css('transform', 'rotate(' + (360 * (ratio - 1)) + 'deg)');
			} else {
				$this.find('.left').css('background-color', param.darkColor).prev().css('transform', 'rotate(' + (360 * (ratio - 0.5)) + 'deg)');
			}
		},
		line : function(param) {
			var $this = param.$this, w = param.w, h = param.h, background = param.background, lineColor = param.frontColor;
			$this.css('position', 'relative').css('background', background);
			var strs = [];
			strs.push('<div style="position:absolute;width:100%;height:100%;top:0;left:0;background:' + lineColor + ';"></div>');
			$this.html(strs.join(''));
			render.lineSet(param);
		},
		lineSet : function(param) {
			var ratio = 100 * param.value / param.total, $this = param.$this;
			$this.find('div').css('width', ratio + '%');
		}
	}

})(window.jQuery || window.Zepto);