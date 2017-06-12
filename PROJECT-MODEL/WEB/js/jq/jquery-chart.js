;
(function($) {

	var PI2 = Math.PI * 2;

	$.fn.chart = function(param, obj) {
		if (typeof param == 'string') {
			return this.filter('canvas').each(function() {
				publicFunc[param] && publicFunc[param]($(this), obj);
			});
		}
		return this.filter('canvas').each(function() {
			$(this).data('chart', init.init(this, $.extend({
				draw : draw
			}, defaultParam, param)));
			render.render(this);
		});
	}

	// 默认参数，传入的参数会将默认参数覆盖掉
	var defaultParam = {
		type : 'area', // 图的类型
		retina : 2,
		data : [ {
			name : '0',
			val : 0
		}, {
			name : '1',
			val : 0
		}, {
			name : '2',
			val : 0
		}, {
			name : '3',
			val : 0
		}, {
			name : '4',
			val : 0
		}, {
			name : '5',
			val : 0
		} ],
		backgroundColor : '#00000000', // 整个图形的背景颜色
		// 绘制饼状图时需要用到的参数：
		antiClockwist : false, // 逆时针
		centerSpace : 2, // 在饼状图中，绘制完成之后由中心点的偏移量
		radius : 0, // 在绘制饼状图时，饼状图中大圆的半径，如果传入该值，则会忽略 radiusRatio 的值。
		radiusRatio : 0.4, // 在饼状图中，半径是长和宽中较短的多少倍
		blankRadius : 0, // 在绘制饼状图时，中间空缺圆形部分的半径的长度，如果传入该值，则会忽略掉 blankRadiusRatio 的值。
		blankRadiusRatio : 0.7, // 在饼状图中，中间空缺圆形部分的半径占整个半径的比例
		circleFont : '2rem 微软雅黑', // 在绘制饼状图时，图上文字的字体
		circleFontColor : '#FFF', // 在绘制饼状图时，图上文字的颜色
		circleLineColor : '#FFF', // 在绘制饼状图时，图上指向文字的线条颜色

		// 绘制坐标轴时需要用到的参数：
		padding : [ 28, 42, 28, 28 ], // 绘制坐标轴时上、右、下、左的边距
		coorColor : '#FFFFFF88', // 坐标轴的颜色
		coorFont : '2rem 微软雅黑', // 坐标轴的字体
		coorFontColor : '#FFFFFF', // 坐标轴上文字的颜色
		coorArea : false, // 是否绘制坐标区域的背景
		coorAreaColor : '#EBF6F9', // 绘制坐标背景时，背景的颜色
		coorGridColor : '#93B8D3', // 绘制坐标网格时，网格的颜色
		xUnit : '', // X 轴的单位
		yUnit : '', // Y 轴的单位
		xTextOffset : 5, // X 轴上文字与 X 轴之间的间隔
		yTextOffset : 0, // Y 轴上文字与 Y 轴之间的间隔
		unitSpace : 2, // xUnit 与 yUnit 被绘制到坐标轴上时，与坐标轴边界的间距
		xType : undefined, // X 轴的绘制方式 area 或 pillar,为空时根据 type 的值而定
		minType : 'fix', // Y 轴的最小值的计算方式
		minTypeVal : 20, // Y 轴的最小值
		maxType : 'fix', // Y 轴的最大值的计算方式
		maxTypeVal : 20, // Y 轴的最大值
		valSteps : 6, // Y 轴被划分为多少个区域
		autoCells : true, // 是否自动计算出 Y 轴被划分为多少个区域
		precision : 0, // 精度：计算结果保留几位小数

		// 绘制面积图时所需要用到的参数：
		areaDot : false, // 绘制面积图时，是否在每个绘制点上用小圆圈
		areaDotColor : '#FFF', // 绘制面积图时，在每个绘制点上的小圆圈的颜色
		areaColor : '#67b138CC', // 绘制面积图时，面积的填充颜色 #29A7D080
		areaFont : '2rem 微软雅黑', // 绘制面积图时，在图上文字的字体
		areaFontColor : '#FFF', // 绘制面积图时，在图上文字的颜色
		areaSmooth : true, // 绘制面积图时，是否采用平滑曲线
		areaFill : true, // 绘制面积图时，是否填充整个面积区域，不填充就相当于折线图

		// 绘制柱状图所需要用到的参数：
		pillarWidth : 0, // 绘制柱状图时，柱状的宽度，不为 0 时会忽略 pillarWidthPercent
		pillarWidthPercent : 0.55, // 绘制柱状图时，柱状宽度是
		pillarFont : '2rem 微软雅黑', // 绘制柱状图时，在柱状图上的文字的字体
		pillarFontColor : '#FFF', // 绘制柱状图时，在柱状图上的文字的颜色

		text : function(i, d) {
			return {
				text : init.fixNum(d.val / this.total * 100, this.precision) + '%',
				font : '2rem 微软雅黑'
			}
		},
		textX : function(i, d) {
			return {
				grid : false,
				label : 0,
				text : d.name
			}
		},
		textY : function(i, v) {
			return {
				grid : false,
				label : 0,
				text : v,
				font : '2rem 微软雅黑'
			}
		}
	}

	// 向外部提供的可直接调用的函数
	var publicFunc = {
		renderArea : function($dom, param) {
			init.initColor(param);
			draw.area($.extend({}, $dom.data('chart'), param));
		},
		loaded : function($dom, param) {
			var data;
			if ($dom && (data = $dom.data('chart')) && data.animFrameId) {
				window.cancelAnimFrame(data.animFrameId);
			}
		}
	}

	var init = {
		// 初始化 Canvas，对数据进行的一些初始化计算
		init : function(dom, param) {
			var $dom = $(dom), g = dom.getContext('2d'), total = 0, min = 1000000000, max = -1000000000;
			// 计算所有值的总和、最小值、最大值
			$.each(param.data, function(i, d) {
				total += parseFloat(d.val);
				min = min < d.val ? min : d.val;
				max = max > d.val ? max : d.val;
			});

			init.initColor(param);

			// 计算宽和高
			$dom.attr('width', param.w = $dom.width() * param.retina);
			$dom.attr('height', param.h = $dom.height() * param.retina);
			for (var i = 0; i < param.padding.length; i++) {
				param.padding[i] *= param.retina;
			}
			g.lineWidth = 2;
			return $.extend(param, {
				g : g,
				total : total || 1,
				max : max,
				min : min
			});
		},
		initColor : function(param) {
			// 对颜色的处理，将#FFFFFF70格式转换成 rbga(255,255,255,0.70);
			$.each(param, function(i, d) {
				if (i.indexOf('Color') != -1) {
					param[i] = draw.colorToRGBA(d);
				}
			});
		},
		initStep : function(param) {
			var xType = param.xType || (param.xType = param.type), xSteps = param.data.length, w = param.w, h = param.h, top = param.padding[0], right = param.padding[1], bottom = param.padding[2], left = param.padding[3];
			var xMin = left, xMax = w - right, yMin = h - bottom, yMax = top;
			// 根据 xType 来计算 X 轴的间隔值
			if (param.xType == 'area') {
				param.xSteps = xSteps - 1;
				param.stepX = (xMax - xMin) / param.xSteps;
				param.xBase = xMin;
			} else if (param.xType == 'pillar') {
				param.xSteps = xSteps + 1;
				param.stepX = (xMax - xMin) / param.xSteps;
				param.xBase = xMin + param.stepX;
			}

			var min = param.min, max = param.max, minType = param.minType, minTypeVal = param.minTypeVal, maxType = param.maxType, maxTypeVal = param.maxTypeVal;
			// 计算 Y 轴的最小值
			if (minType == 'min') {
				min = min * minTypeVal;
			} else if (minType == 'num') {
				min = minTypeVal;
			} else if (minType == 'fix') {
				min = Math.floor(min / minTypeVal) * minTypeVal;
			}
			param.minVal = min;

			// 计算 Y 轴的最大值
			if (maxType == 'max') {
				max = max * maxTypeVal;
			} else if (maxType == 'num') {
				max = maxTypeVal;
			} else if (maxType == 'fix') {
				max = (Math.ceil(max / maxTypeVal) || 1) * maxTypeVal;
			}
			param.maxVal = max;

			// 自动计算 Y 轴的间隔值
			if (param.autoCells) {
				if (minType == 'fix' && maxType == 'fix' && minTypeVal == maxTypeVal) {
					param.valSteps = (param.maxVal - param.minVal) / maxTypeVal;
					param.stepVal = maxTypeVal;
				} else {
					param.stepVal = (param.maxVal - param.minVal) / param.valSteps;
				}
			} else {
				param.stepVal = (param.maxVal - param.minVal) / param.valSteps;
			}

			param.ySteps = param.valSteps;
			param.stepY = (yMin - yMax) / param.ySteps;

		},
		initPillarWidth : function(param) {
			param.pillarWidth = param.pillarWidth || (param.stepX * param.pillarWidthPercent / 2);
		},
		fixNum : function(num, precision) {
			return num.toFixed(precision);
		}

	}

	// 在面板上绘制单一图行
	var draw = {
		colors : [ '#9D4A4B', '#5E7F99', '#98B3BD', '#A5AAAC', '#787F89', '#6F83A5', '#FFFFCC', '#CCFFFF', '#FFCCCC', '#99CCCC', '#CCFFCC', '#CCCCFF', '#99CCFF', '#FFCC99', '#336699', '#CCFF99', '#CCCCFF', '#99CC33', '#0099CC', '#CCCCCC' ],
		randomColor : function(data) {
			var tempColors = $.extend([], this.colors);
			$.each(data, function(i, d) {
				for (var i = 0; i < tempColors.length; i++) {
					if (tempColors[i] == d.color) {
						tempColors.splice(i, 1);
						break;
					}
				}
			});
			return tempColors.length ? tempColors[Math.floor(Math.random() * tempColors.length)] : 'black';
		},
		colorToRGBA : function(color) {
			var colors;
			if (color && color.length == 9 && (colors = /^#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$/.exec(color))) {
				return [ 'rgba(', parseInt(colors[1], 16), ',', parseInt(colors[2], 16), ',', parseInt(colors[3], 16), ',', parseInt(colors[4], 16) / 256, ')' ].join('');
			} else {
				return color;
			}
		},
		colorToGradient : function(color, g, x1, y1, x2, y2) {
			var colos;
			if (color && (colors = /^(#[0-9a-fA-F]{6,8})-(#[0-9a-fA-F]{6,8})/.exec(color))) {
				var lGrd = g.createLinearGradient(x1, y1, x2, y2);
				lGrd.addColorStop(0, draw.colorToRGBA(colors[1]));
				lGrd.addColorStop(1, draw.colorToRGBA(colors[2]));
				return lGrd;
			} else {
				return draw.colorToRGBA(color);
			}
		},
		background : function(param) {
			var g = param.g;
			g.fillStyle = param.backgroundColor;
			g.fillRect(0, 0, param.w, param.h);
		},
		arc : function(g, x, y, r, blankRadius, fromAngle, toAngle, antiClockwist, fillStyle) {
			g.fillStyle = fillStyle;
			g.beginPath();
			g.arc(x, y, r, fromAngle, toAngle, antiClockwist);
			g.arc(x, y, blankRadius, toAngle, fromAngle, !antiClockwist);
			g.closePath();
			g.fill();
		},
		textInCircle : function(g, x, y, r, angle, text, fontSize, fontColor, lineColor) {
			var x1 = x + r * Math.cos(angle);
			var y1 = y + r * Math.sin(angle);
			var x2 = x1 + 10 * Math.cos(angle);
			var y2 = y1 + 10 * Math.sin(angle);
			draw.line(g, [ x1, y1, x2, y2 ], lineColor);
			draw.fillText(g, text, x2, y2, x2 > x1 ? 'left' : 'right', y2 > y1 ? 'top' : 'bottom', fontColor, fontSize);
		},
		line : function(g, coors, style, fill, loop) {
			var i = 1, len = coors.length >> 1;
			g.beginPath();
			g.moveTo(coors[0], coors[1]);
			for (; i < len; i++) {
				g.lineTo(coors[i << 1], coors[(i << 1) + 1]);
			}
			style && (fill ? (g.fillStyle = style) : (g.strokeStyle = style));
			if (loop) {
				g.closePath();
				fill ? g.fill() : g.stroke();
			} else {
				fill ? g.fill() : g.stroke();
				g.closePath();
			}
		},
		getRefPoint : function(x1, y1, x2, y2, x3, y3, l) {
			var da = Math.sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1)), va = [ (x2 - x1) / da, (y2 - y1) / da ];
			var db = Math.sqrt((x3 - x2) * (x3 - x2) + (y3 - y2) * (y3 - y2)), vb = [ (x3 - x2) / db, (y3 - y2) / db ];
			var vt = [ va[0] + vb[0], va[1] + vb[1] ], dt = Math.sqrt(vt[0] * vt[0] + vt[1] * vt[1]);
			vt = [ vt[0] / dt, vt[1] / dt ];

			return {
				x1 : x2 + vt[0] * l,
				y1 : y2 + vt[1] * l,
				x2 : x2 - vt[0] * l,
				y2 : y2 - vt[1] * l
			}
		},
		bessel : function(g, coors, offset, yMin) {
			var i = 0, len = coors.length, x0, y0, x1, y1, x2, y2, x3, y2, sx, sy, oldPoints, refPoints;
			g.lineTo(coors[0], coors[1]);
			for (; i < len; i += 2) {
				if (i >= len - 4) {
					refPoints = draw.getRefPoint(coors[i], coors[i + 1], coors[i + 2], coors[i + 3], coors[i + 2] + 1, coors[i + 3], offset);
				} else {
					refPoints = draw.getRefPoint(coors[i], coors[i + 1], coors[i + 2], coors[i + 3], coors[i + 4], coors[i + 5], offset);
				}
				refPoints.y1 = Math.min(refPoints.y1, yMin);
				refPoints.y2 = Math.min(refPoints.y2, yMin);
				if (!oldPoints && refPoints.x2 && refPoints.y2) {
					g.quadraticCurveTo(refPoints.x2, refPoints.y2, coors[i + 2], coors[i + 3]);
				} else if (oldPoints && oldPoints.x1 && oldPoints.y1 && refPoints.x2 && refPoints.y2) {
					g.bezierCurveTo(oldPoints.x1, oldPoints.y1, refPoints.x2, refPoints.y2, coors[i + 2], coors[i + 3]);
				} else if (oldPoints && oldPoints.x2 && oldPoints.y2) {
					g.quadraticCurveTo(oldPoints.x1, oldPoints.y1, coors[i + 2], coors[i + 3]);
				}
				oldPoints = refPoints;
			}
		},
		curve : function(g, coors, style, fill, loop, offset, yMin) {
			var len = coors.length, lastX = coors[len - 2], lastY = coors[len - 1];
			g.beginPath();
			g.moveTo(coors[0], coors[1]);
			draw.bessel(g, coors.splice(2, coors.length - 4), offset, yMin);
			g.lineTo(lastX, lastY);
			style && (fill ? (g.fillStyle = style) : (g.strokeStyle = style));
			if (loop) {
				g.closePath();
				fill ? g.fill() : g.stroke();
			} else {
				fill ? g.fill() : g.stroke();
				g.closePath();
			}
		},
		fillText : function(g, text, x, y, textAlign, textBaseline, fillStyle, font) {
			g.textAlign = textAlign;
			g.textBaseline = textBaseline;
			fillStyle && (g.fillStyle = fillStyle);
			font && (g.font = font);
			g.fillText(text, x, y);
		},
		coordinate : function(param) {
			var g = param.g, i, top = param.padding[0], right = param.padding[1], bottom = param.padding[2], left = param.padding[3], w = param.w, h = param.h;
			var xMin = left, xMax = w - right, yMin = h - bottom, yMax = top;

			// 绘制背景区域
			if (param.coorArea) {
				g.fillStyle = param.coorAreaColor;
				g.fillRect(xMin, yMax, xMax - xMin, yMin - yMax);
			}

			// 绘制网格、刻度以及刻度上的文字
			var xSteps = param.xSteps, stepX = param.stepX, ySteps = param.ySteps, stepY = param.stepY, stepVal = param.stepVal, xBase = param.xBase, xTextOffset = param.xTextOffset, yTextOffset = param.yTextOffset, pillarWidth = param.pillarWidth, tempObject;
			for (g.strokeStyle = param.coorColor, i = 0; i <= xSteps; i++) {
				if (tempObject = param.textX(i, param.data[i])) {
					tempObject.grid && this.line(g, [ xBase + stepX * i, yMin, xBase + stepX * i, yMax ]);
					tempObject.label && this.line(g, [ xBase + stepX * i, yMin, xBase + stepX * i, yMin - tempObject.label ]);
					tempObject.text && param.data[i] && this.fillText(g, tempObject.text, xBase + stepX * i, yMin + yTextOffset, 'center', 'top', param.coorFontColor, tempObject.font || param.coorFont);
				}
			}
			for (g.strokeStyle = param.coorColor, i = 0; i <= ySteps; i++) {
				if (tempObject = param.textY(i, init.fixNum(param.minVal + stepVal * i, param.precision))) {
					tempObject.grid && this.line(g, [ xMin, yMin - stepY * i, xMax, yMin - stepY * i ]);
					tempObject.label && this.line(g, [ xMin, yMin - stepY * i, xMin + tempObject.label, yMin - stepY * i ]);
					tempObject.text && this.fillText(g, tempObject.text, xMin - xTextOffset, yMin - stepY * i, 'right', 'middle', param.coorFontColor, tempObject.font || param.coorFont);
				}
			}

			// 绘制 X 轴
			g.lineWidth = 4;
			this.line(g, [ xMin - (pillarWidth ? pillarWidth : 1), yMin + 1, xMax + (pillarWidth ? pillarWidth : 1), yMin + 1 ], param.coorColor);
			// 绘制 X 轴、Y 轴上的单位
			this.fillText(g, param.xUnit, xMax + param.unitSpace, yMin, 'left', 'bottom', param.coorFontColor, param.coorFont);
			this.fillText(g, param.yUnit, xMin, yMax - param.unitSpace, 'left', 'bottom', param.coorFontColor, param.coorFont);
		},
		circleChart : function(param) {
			var fromPercent = 0, toPercent = 0, fromAngle, toAngle, directionAngle, centerX, centerY, g = param.g, w = param.w, h = param.h, r = param.radius || (Math.min(w, h) * param.radiusRatio), centerSpace = param.centerSpace || 0, antiClockwist = param.antiClockwist, blankRadius = param.blankRadius || (r * param.blankRadiusRatio), temp;
			$.each(param.data, function(i, d) {
				fromAngle = PI2 * fromPercent;
				toAngle = PI2 * (fromPercent += (antiClockwist ? -d.val : d.val) / param.total);
				directionAngle = fromAngle + ((toAngle - fromAngle) / 2);
				centerX = w / 2 + centerSpace * Math.cos(directionAngle);
				centerY = h / 2 + centerSpace * Math.sin(directionAngle);
				draw.arc(g, centerX, centerY, r, blankRadius, fromAngle, toAngle, param.antiClockwist, d.color || (d.color = draw.randomColor(param.data)));
				(temp = param.text(i, d)) && draw.textInCircle(g, centerX, centerY, r, directionAngle, temp.text, temp.font, param.circleFontColor, param.circleLineColor);
			});
		},
		area : function(param) {
			var g = param.g, x = 0, y = 0, coors = [], top = param.padding[0], right = param.padding[1], bottom = param.padding[2], left = param.padding[3], w = param.w, h = param.h, minVal = param.minVal, maxVal = param.maxVal;
			var xSteps = param.xSteps, stepX = param.stepX, ySteps = param.ySteps, stepY = param.stepY, stepVal = param.stepVal, xBase = param.xBase, xMin = left, xMax = w - right, yMin = h - bottom, yMax = top, data = param.data, temp;

			coors.push(xBase, yMin);
			$.each(data, function(i, d) {
				d.x = x = xBase + stepX * i;
				d.y = yMin - (yMin - yMax) * d.val / (maxVal - minVal);
				coors.push(d.x, d.y);
			});
			coors.push(x, yMin);
			param.areaColor = draw.colorToGradient(param.areaColor, g, xMin, yMax, xMin, yMin);
			draw[param.areaSmooth ? 'curve' : 'line'](g, coors, param.areaColor, param.areaFill, false, stepX / 2, yMin);

			$.each(data, function(i, d) {
				param.areaDot && draw.arc(g, d.x, d.y, 5, 2, 0, PI2, false, param.areaDotColor);
				if ((temp = param.text(i, d))) {
					var fontSize = parseInt(temp.font || param.areaFont);
					var textW = fontSize * temp.text.length;
					var textH = fontSize;
					if (d.x - textW / 2 > xMin && d.x + textW / 2 < xMax) {
						x = d.x, y = d.y;
					} else if (d.x - textW / 2 <= xMin) {
						x = xMin + textW / 2 - fontSize, y = d.y;
					} else {
						x = xMax - textW / 2 + fontSize, y = d.y;
					}
					temp.border && draw.line(g, [ d.x, d.y, d.x - fontSize / 2, d.y - fontSize / 2 - 2, d.x + fontSize / 2, d.y - fontSize / 2 - 2 ], undefined, true);
					temp.border && g.fillRect(x - textW / 2, y - fontSize * 2, textW, fontSize * 1.5);
					temp.text && draw.fillText(g, temp.text, x, y - fontSize * 0.5, 'center', 'bottom', param.areaFontColor, temp.font);
				}
			});
		},
		pillar : function(param) {
			var g = param.g, x = 0, y = 0, coors = [], top = param.padding[0], right = param.padding[1], bottom = param.padding[2], left = param.padding[3], w = param.w, h = param.h, minVal = param.minVal, maxVal = param.maxVal;
			var xSteps = param.xSteps, stepX = param.stepX, ySteps = param.ySteps, stepY = param.stepY, stepVal = param.stepVal, xBase = param.xBase, xMin = left, xMax = w - right, yMin = h - bottom, yMax = top, pillarWidth = param.pillarWidth, borderRadius = pillarWidth / 3, temp;
			$.each(param.data, function(i, d) {
				coors = [];
				x = xBase + stepX * i;
				y = yMin - (yMin - yMax) / (maxVal - minVal) * (d.val - minVal);
				if ($.isFunction(d.draw) && d.draw(param, x, y) === true) {
					d.draw(param);
					return true;
				}
				coors.push(x - pillarWidth, y + borderRadius);
				coors.push(x - pillarWidth + borderRadius, y + borderRadius);
				coors.push(x - pillarWidth + borderRadius, y);
				coors.push(x + pillarWidth - borderRadius, y);
				coors.push(x + pillarWidth - borderRadius, y + borderRadius);
				coors.push(x + pillarWidth, y + borderRadius);
				coors.push(x + pillarWidth, yMin);
				coors.push(x - pillarWidth, yMin);
				draw.line(g, coors, (d.color = draw.colorToGradient(d.color, g, x, y, x, yMin)) || (d.color = draw.randomColor(param.data)), true);
				draw.arc(g, x - pillarWidth + borderRadius + 1, y + borderRadius + 1, borderRadius + 1, 0, Math.PI, Math.PI + Math.PI / 2, false, d.color);
				draw.arc(g, x + pillarWidth - borderRadius - 1, y + borderRadius + 1, borderRadius + 1, 0, -Math.PI / 2, 0, false, d.color);
				(temp = param.text(i, d)) && draw.fillText(g, temp.text, x, y, 'center', 'bottom', param.pillarFontColor, temp.font);
				// draw.fillText(g, param.text(d), x, yMin - (yMin - y) / 2, 'center', 'middle', param.pillarFontColor, param.pillarFont);
			});
		},
		loading : function(param) {
			var g = param.g, w = param.w, h = param.h;
			var angle = 0, r = Math.min(100, w), r2 = r / 2, color = "#000000";
			var startTime = new Date().getTime(), endTime;
			(function refresh() {
				endTime = new Date().getTime();
				angle = (endTime - startTime) / 180;
				g.save();
				g.translate(w / 2, h / 2);
				g.rotate(angle);
				g.clearRect(-r2, -r2, r, r);
				for (var i = 0; i < 15; i++) {
					g.rotate(0.4);
					g.fillStyle = draw.colorToRGBA(color + ((i + 1) * parseInt(11, 16)).toString(16));
					g.fillRect(-r2 * 0.05, -r2, r2 * 0.1, r2 * 0.6);
				}
				g.restore();
				param.animFrameId = window.requestAnimFrame(refresh);
				// console.info(param.animFrameId);
			})();
		}
	}

	var render = {
		// 在 Canvas 中绘制整个图行
		render : function(dom) {
			var param = $(dom).data('chart');
			render[param.type](param);
		},
		circle : function(param) {
			draw.background(param);
			draw.circleChart(param);
		},
		area : function(param) {
			init.initStep(param);
			draw.background(param);
			draw.coordinate(param);
			draw.area(param);
		},
		pillar : function(param) {
			init.initStep(param);
			init.initPillarWidth(param);
			draw.background(param);
			draw.coordinate(param);
			draw.pillar(param);
		},
		loading : function(param) {
			init.initStep(param);
			draw.background(param);
			draw.coordinate(param);
			draw.loading(param);
		}
	}

})(window.jQuery || window.Zepto);