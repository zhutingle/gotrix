;
(function() {

	var plugin = 'slide';

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
				$(this).data(plugin, new Slide(this, param));
			});
		}
	}

	var defaultParam = {
		data : []
	}

	function Slide(dom, param) {
		$.extend(this, defaultParam, param);
		var _this = this;
		_this.$dom = $(dom).css({
			'position' : 'relative',
			'overflow' : 'hidden',
			'width' : '100%'
		});

		var strs = [];
		strs.push('<ul style="width:' + _this.data.length + '00%">');
		strs.push(('<li style="float:left;width:' + (100 / _this.data.length) + '%;"><img src="*{img}" style="display:block;height:10rem;margin:0 auto;" /></li>').from(_this.data))
		strs.push('</ul>');
		strs.push('<div style="clear:both;"></div>');
		strs.push('<div style="position: absolute;width: 100%;text-align:center;bottom: 0rem;">');
		strs.push('<span style="display:inline-block;width: 0.2rem;height: 0.2rem;border: 0.05rem solid white;margin: 0.05rem;border-radius: 100rem;"></span>'.from(_this.data));
		strs.push('</div>')
		_this.$dom.html(strs.join(''));

		_this.move();
		_this.$dom.delegate('img', 'click', function() {
			var d = _this.data[$(this).index()];
			if (d.url) {
				window.location = d.url;
			}
		});
	}

	$.extend(Slide.prototype, {
		moveTo : function(index) {
			this.index = index;
			this.$dom.find('span:eq(' + index + ')').css('background', '#FF0000').siblings().css('background', '#9E8780');
			this.$dom.find('ul').stop().animate({
				marginLeft : '-' + (100 * index) + '%'
			}, 1000);
		},
		moveToNext : function() {
			this.moveTo(this.index + 1 < this.data.length ? this.index + 1 : 0);
		},
		move : function() {
			var _this = this;
			_this.moveTo(0);
			var interval = setInterval(function() {
				_this.moveToNext();
			}, 3000);
		}
	});

})(window.jQuery || window.Zepto);