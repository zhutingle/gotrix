(function($) {

	// -------------------------------------------------------
	// $('.page-current').removeClass('page-current').addClass('hide');
	// $('.page-7-1').removeClass('hide').addClass('page-current');
	// -------------------------------------------------------

	var $curPage, now;
	var towards = {
		up : 1,
		right : 2,
		down : 3,
		left : 4
	};
	var noSwipe = false;
	var onlyLRSwipe = false;
	var isAnimating = false;

	$(function() {

		$curPage = $('.page-current');
		if ($curPage[0]) {
			var curPages = /page-(\d+)-(\d+)/.exec($curPage[0].className);
			now = {
				row : parseInt(curPages[1]),
				col : parseInt(curPages[2])
			};
		} else {
			now = {
				row : 1,
				col : 1
			};
		}

		var waitingHeight = $(".waiting").height();
		var waitingWidth = $(".waiting").width();
		window.repairTrans = function() {
			var height = $(window).height() > 2 * waitingHeight ? waitingHeight : $(window).height();
			var s = height / parseInt($('.wrap').css('height'));
			
			var width = $(window).height() > 2 * waitingHeight ? waitingWidth : $(window).width();
			var ws = width / parseInt($('.wrap').css('width'));
			if(ws > s * 1.5) {
				ws = s;
			}
			var transform = 'scale(' + ws + ',' + s + ')';
			$('.wrap').css({
				'-webkit-transform' : transform,
				'-moz-transform' : transform,
				'-o-transform' : transform
			});
		}

		var $imgs = $("img");
		var loadedInterval = setInterval(function() {
			var allComplete = false;
			$imgs.each(function() {
				return allComplete = this.complete;
			});
			allComplete = window.loadComplete ? window.loadComplete() : true;
			if ($imgs.length == 0 || allComplete) {
				$(window).resize(repairTrans);
				repairTrans();
				$('.swiper-wrapper,.page-current').removeClass('hide');
				$(".waiting").remove();
				clearInterval(loadedInterval);
				if (typeof $ready == 'function') {
					$ready();
				}
			}
		}, 100);
	});

	var oldTouch;

	function getTouch(e) {
		function isPrimaryTouch(event) {
			return (event.pointerType == 'touch' || event.pointerType == event.MSPOINTER_TYPE_TOUCH) && event.isPrimary;
		}
		function isPointerEventType(e, type) {
			return (e.type == 'pointer' + type || e.type.toLowerCase() == 'mspointer' + type);
		}
		if ((_isPointerType = isPointerEventType(e, 'move')) && !isPrimaryTouch(e)) {
			return;
		}
		return _isPointerType ? e : e.touches[0];
	}

	document.addEventListener('touchstart', function(e) {
		oldTouch = getTouch(e);
	});

	document.addEventListener('touchmove', function(e) {
		if (noSwipe) {
			return;
		}
		if (onlyLRSwipe) {
			if (Math.abs(oldTouch.pageX - getTouch(e).pageX) > 5) {
				e.preventDefault();
			}
			return;
		}
		e.preventDefault();
	}, false);

	bindMove("swipeUp", "icon-up", function() {
		if (isAnimating)
			return;
		pageMove(towards.up, 1, 0);
	});
	bindMove("swipeDown", "icon-down", function() {
		if (isAnimating)
			return;
		pageMove(towards.down, -1, 0);
	});
	bindMove("swipeLeft", "icon-left,.back-left", function() {
		if (isAnimating)
			return;
		pageMove(towards.left, 0, 1);
	});
	bindMove("swipeRight", "icon-right,.back-right", function() {
		if (isAnimating)
			return;
		pageMove(towards.right, 0, -1);
	});

	function bindMove(action, icon, callback) {
		$(document)[action](function() {
			if (noSwipe) {
				return;
			}
			if (onlyLRSwipe && /swipe(Up)|(Down)/.test(action)) {
				return;
			}
			callback();
		});
		$(document).delegate('.' + icon, 'tap', callback);
		// $('.' + icon).on('tap', callback);
	}

	function pageMove(tw, rowStep, colStep) {
		var lastPage = ".page-" + now.row + "-" + now.col;
		var row = now.row, col = now.col;
		var nowPage;
		var findNowPage = false;
		for (var i = 0; i < 10; i++) {
			if (rowStep != 0) { // 如果row有变化，则强制将col列变为1。
				colStep = 1 - now.col;
				nowPage = ".page-" + (row += rowStep) + "-1";
			} else {
				nowPage = ".page-" + (row += rowStep) + "-" + (col += colStep);
			}
			window.onPageMove && window.onPageMove(tw, rowStep, colStep, nowPage);
			if ($(nowPage)[0]) {
				findNowPage = true;
				break;
			}
		}
		if (!findNowPage) {
			return;
		}
		now.row = row;
		now.col = col;

		// now.row += rowStep;
		// now.col += colStep;
		// if (now.row == 7 && (now.col == 2 || now.col == 3)) {
		// noSwipe = true;
		// } else {
		// $(window).scrollTop(0);
		// noSwipe = false;
		// }

		switch (tw) {
		case towards.up:
			outClass = 'pt-page-moveToTop';
			inClass = 'pt-page-moveFromBottom';
			break;
		case towards.right:
			outClass = 'pt-page-moveToRight';
			inClass = 'pt-page-moveFromLeft';
			break;
		case towards.down:
			outClass = 'pt-page-moveToBottom';
			inClass = 'pt-page-moveFromTop';
			break;
		case towards.left:
			outClass = 'pt-page-moveToLeft';
			inClass = 'pt-page-moveFromRight';
			break;
		}
		isAnimating = true;
		$(nowPage).removeClass("hide");

		$(lastPage).addClass(outClass);
		$(nowPage).addClass(inClass);

		setTimeout(function() {

			$(lastPage).removeClass('page-current').removeClass(outClass).addClass("hide").find("img").addClass("hide");
			$(nowPage).addClass('page-current').removeClass(inClass).find("img").removeClass("hide");

			setTimeout(function() {
				isAnimating = false;
			}, 100);
		}, 400);
	}

})(window.Zepto);