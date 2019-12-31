/* ========================================================================
 * bootstrap-switch - v3.3.2
 * http://www.bootstrap-switch.org
 * ========================================================================
 * Copyright 2012-2013 Mattia Larentis
 *
 * ========================================================================
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ========================================================================
 */

(function() {
    var slice = [].slice;

    (function($, window) {
        "use strict";
        var BootstrapSwitch;
        BootstrapSwitch = (function() {
            function BootstrapSwitch(element, options) {
                if (options == null) {
                    options = {};
                }
                this.$element = $(element);
                this.options = $.extend({}, $.fn.bootstrapSwitch.defaults, {
                    state: this.$element.is(":checked"),
                    size: this.$element.data("size"),
                    animate: this.$element.data("animate"),
                    disabled: this.$element.is(":disabled"),
                    readonly: this.$element.is("[readonly]"),
                    indeterminate: this.$element.data("indeterminate"),
                    inverse: this.$element.data("inverse"),
                    radioAllOff: this.$element.data("radio-all-off"),
                    onColor: this.$element.data("on-color"),
                    offColor: this.$element.data("off-color"),
                    onText: this.$element.data("on-text"),
                    offText: this.$element.data("off-text"),
                    labelText: this.$element.data("label-text"),
                    handleWidth: this.$element.data("handle-width"),
                    labelWidth: this.$element.data("label-width"),
                    baseClass: this.$element.data("base-class"),
                    wrapperClass: this.$element.data("wrapper-class")
                }, options);
                this.prevOptions = {};
                this.$wrapper = $("<div>", {
                    "class": (function(_this) {
                        return function() {
                            var classes;
                            classes = ["" + _this.options.baseClass].concat(_this._getClasses(_this.options.wrapperClass));
                            classes.push(_this.options.state ? _this.options.baseClass + "-on" : _this.options.baseClass + "-off");
                            if (_this.options.size != null) {
                                classes.push(_this.options.baseClass + "-" + _this.options.size);
                            }
                            if (_this.options.disabled) {
                                classes.push(_this.options.baseClass + "-disabled");
                            }
                            if (_this.options.readonly) {
                                classes.push(_this.options.baseClass + "-readonly");
                            }
                            if (_this.options.indeterminate) {
                                classes.push(_this.options.baseClass + "-indeterminate");
                            }
                            if (_this.options.inverse) {
                                classes.push(_this.options.baseClass + "-inverse");
                            }
                            if (_this.$element.attr("id")) {
                                classes.push(_this.options.baseClass + "-id-" + (_this.$element.attr("id")));
                            }
                            return classes.join(" ");
                        };
                    })(this)()
                });
                this.$container = $("<div>", {
                    "class": this.options.baseClass + "-container"
                });
                this.$on = $("<span>", {
                    html: this.options.onText,
                    "class": this.options.baseClass + "-handle-on " + this.options.baseClass + "-" + this.options.onColor
                });
                this.$off = $("<span>", {
                    html: this.options.offText,
                    "class": this.options.baseClass + "-handle-off " + this.options.baseClass + "-" + this.options.offColor
                });
                this.$label = $("<span>", {
                    html: this.options.labelText,
                    "class": this.options.baseClass + "-label"
                });
                this.$element.on("init.bootstrapSwitch", (function(_this) {
                    return function() {
                        return _this.options.onInit.apply(element, arguments);
                    };
                })(this));
                this.$element.on("switchChange.bootstrapSwitch", (function(_this) {
                    return function(e) {
                        if (false === _this.options.onSwitchChange.apply(element, arguments)) {
                            if (_this.$element.is(":radio")) {
                                return $("[name='" + (_this.$element.attr('name')) + "']").trigger("previousState.bootstrapSwitch", true);
                            } else {
                                return _this.$element.trigger("previousState.bootstrapSwitch", true);
                            }
                        }
                    };
                })(this));
                this.$container = this.$element.wrap(this.$container).parent();
                this.$wrapper = this.$container.wrap(this.$wrapper).parent();
                this.$element.before(this.options.inverse ? this.$off : this.$on).before(this.$label).before(this.options.inverse ? this.$on : this.$off);
                if (this.options.indeterminate) {
                    this.$element.prop("indeterminate", true);
                }
                this._init();
                this._elementHandlers();
                this._handleHandlers();
                this._labelHandlers();
                this._formHandler();
                this._externalLabelHandler();
                this.$element.trigger("init.bootstrapSwitch", this.options.state);
            }

            BootstrapSwitch.prototype._constructor = BootstrapSwitch;

            BootstrapSwitch.prototype.setPrevOptions = function() {
                return this.prevOptions = $.extend(true, {}, this.options);
            };

            BootstrapSwitch.prototype.state = function(value, skip) {
                if (typeof value === "undefined") {
                    return this.options.state;
                }
                if (this.options.disabled || this.options.readonly) {
                    return this.$element;
                }
                if (this.options.state && !this.options.radioAllOff && this.$element.is(":radio")) {
                    return this.$element;
                }
                if (this.$element.is(":radio")) {
                    $("[name='" + (this.$element.attr('name')) + "']").trigger("setPreviousOptions.bootstrapSwitch");
                } else {
                    this.$element.trigger("setPreviousOptions.bootstrapSwitch");
                }
                if (this.options.indeterminate) {
                    this.indeterminate(false);
                }
                value = !!value;
                this.$element.prop("checked", value).trigger("change.bootstrapSwitch", skip);
                return this.$element;
            };

            BootstrapSwitch.prototype.toggleState = function(skip) {
                if (this.options.disabled || this.options.readonly) {
                    return this.$element;
                }
                if (this.options.indeterminate) {
                    this.indeterminate(false);
                    return this.state(true);
                } else {
                    return this.$element.prop("checked", !this.options.state).trigger("change.bootstrapSwitch", skip);
                }
            };

            BootstrapSwitch.prototype.size = function(value) {
                if (typeof value === "undefined") {
                    return this.options.size;
                }
                if (this.options.size != null) {
                    this.$wrapper.removeClass(this.options.baseClass + "-" + this.options.size);
                }
                if (value) {
                    this.$wrapper.addClass(this.options.baseClass + "-" + value);
                }
                this._width();
                this._containerPosition();
                this.options.size = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.animate = function(value) {
                if (typeof value === "undefined") {
                    return this.options.animate;
                }
                value = !!value;
                if (value === this.options.animate) {
                    return this.$element;
                }
                return this.toggleAnimate();
            };

            BootstrapSwitch.prototype.toggleAnimate = function() {
                this.options.animate = !this.options.animate;
                this.$wrapper.toggleClass(this.options.baseClass + "-animate");
                return this.$element;
            };

            BootstrapSwitch.prototype.disabled = function(value) {
                if (typeof value === "undefined") {
                    return this.options.disabled;
                }
                value = !!value;
                if (value === this.options.disabled) {
                    return this.$element;
                }
                return this.toggleDisabled();
            };

            BootstrapSwitch.prototype.toggleDisabled = function() {
                this.options.disabled = !this.options.disabled;
                this.$element.prop("disabled", this.options.disabled);
                this.$wrapper.toggleClass(this.options.baseClass + "-disabled");
                return this.$element;
            };

            BootstrapSwitch.prototype.readonly = function(value) {
                if (typeof value === "undefined") {
                    return this.options.readonly;
                }
                value = !!value;
                if (value === this.options.readonly) {
                    return this.$element;
                }
                return this.toggleReadonly();
            };

            BootstrapSwitch.prototype.toggleReadonly = function() {
                this.options.readonly = !this.options.readonly;
                this.$element.prop("readonly", this.options.readonly);
                this.$wrapper.toggleClass(this.options.baseClass + "-readonly");
                return this.$element;
            };

            BootstrapSwitch.prototype.indeterminate = function(value) {
                if (typeof value === "undefined") {
                    return this.options.indeterminate;
                }
                value = !!value;
                if (value === this.options.indeterminate) {
                    return this.$element;
                }
                return this.toggleIndeterminate();
            };

            BootstrapSwitch.prototype.toggleIndeterminate = function() {
                this.options.indeterminate = !this.options.indeterminate;
                this.$element.prop("indeterminate", this.options.indeterminate);
                this.$wrapper.toggleClass(this.options.baseClass + "-indeterminate");
                this._containerPosition();
                return this.$element;
            };

            BootstrapSwitch.prototype.inverse = function(value) {
                if (typeof value === "undefined") {
                    return this.options.inverse;
                }
                value = !!value;
                if (value === this.options.inverse) {
                    return this.$element;
                }
                return this.toggleInverse();
            };

            BootstrapSwitch.prototype.toggleInverse = function() {
                var $off, $on;
                this.$wrapper.toggleClass(this.options.baseClass + "-inverse");
                $on = this.$on.clone(true);
                $off = this.$off.clone(true);
                this.$on.replaceWith($off);
                this.$off.replaceWith($on);
                this.$on = $off;
                this.$off = $on;
                this.options.inverse = !this.options.inverse;
                return this.$element;
            };

            BootstrapSwitch.prototype.onColor = function(value) {
                var color;
                color = this.options.onColor;
                if (typeof value === "undefined") {
                    return color;
                }
                if (color != null) {
                    this.$on.removeClass(this.options.baseClass + "-" + color);
                }
                this.$on.addClass(this.options.baseClass + "-" + value);
                this.options.onColor = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.offColor = function(value) {
                var color;
                color = this.options.offColor;
                if (typeof value === "undefined") {
                    return color;
                }
                if (color != null) {
                    this.$off.removeClass(this.options.baseClass + "-" + color);
                }
                this.$off.addClass(this.options.baseClass + "-" + value);
                this.options.offColor = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.onText = function(value) {
                if (typeof value === "undefined") {
                    return this.options.onText;
                }
                this.$on.html(value);
                this._width();
                this._containerPosition();
                this.options.onText = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.offText = function(value) {
                if (typeof value === "undefined") {
                    return this.options.offText;
                }
                this.$off.html(value);
                this._width();
                this._containerPosition();
                this.options.offText = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.labelText = function(value) {
                if (typeof value === "undefined") {
                    return this.options.labelText;
                }
                this.$label.html(value);
                this._width();
                this.options.labelText = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.handleWidth = function(value) {
                if (typeof value === "undefined") {
                    return this.options.handleWidth;
                }
                this.options.handleWidth = value;
                this._width();
                this._containerPosition();
                return this.$element;
            };

            BootstrapSwitch.prototype.labelWidth = function(value) {
                if (typeof value === "undefined") {
                    return this.options.labelWidth;
                }
                this.options.labelWidth = value;
                this._width();
                this._containerPosition();
                return this.$element;
            };

            BootstrapSwitch.prototype.baseClass = function(value) {
                return this.options.baseClass;
            };

            BootstrapSwitch.prototype.wrapperClass = function(value) {
                if (typeof value === "undefined") {
                    return this.options.wrapperClass;
                }
                if (!value) {
                    value = $.fn.bootstrapSwitch.defaults.wrapperClass;
                }
                this.$wrapper.removeClass(this._getClasses(this.options.wrapperClass).join(" "));
                this.$wrapper.addClass(this._getClasses(value).join(" "));
                this.options.wrapperClass = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.radioAllOff = function(value) {
                if (typeof value === "undefined") {
                    return this.options.radioAllOff;
                }
                value = !!value;
                if (value === this.options.radioAllOff) {
                    return this.$element;
                }
                this.options.radioAllOff = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.onInit = function(value) {
                if (typeof value === "undefined") {
                    return this.options.onInit;
                }
                if (!value) {
                    value = $.fn.bootstrapSwitch.defaults.onInit;
                }
                this.options.onInit = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.onSwitchChange = function(value) {
                if (typeof value === "undefined") {
                    return this.options.onSwitchChange;
                }
                if (!value) {
                    value = $.fn.bootstrapSwitch.defaults.onSwitchChange;
                }
                this.options.onSwitchChange = value;
                return this.$element;
            };

            BootstrapSwitch.prototype.destroy = function() {
                var $form;
                $form = this.$element.closest("form");
                if ($form.length) {
                    $form.off("reset.bootstrapSwitch").removeData("bootstrap-switch");
                }
                this.$container.children().not(this.$element).remove();
                this.$element.unwrap().unwrap().off(".bootstrapSwitch").removeData("bootstrap-switch");
                return this.$element;
            };

            BootstrapSwitch.prototype._width = function() {
                var $handles, handleWidth;
                $handles = this.$on.add(this.$off);
                $handles.add(this.$label).css("width", "");
                handleWidth = this.options.handleWidth === "auto" ? Math.max(this.$on.width(), this.$off.width()) : this.options.handleWidth;
                $handles.width(handleWidth);
                this.$label.width((function(_this) {
                    return function(index, width) {
                        if (_this.options.labelWidth !== "auto") {
                            return _this.options.labelWidth;
                        }
                        if (width < handleWidth) {
                            return handleWidth;
                        } else {
                            return width;
                        }
                    };
                })(this));
                this._handleWidth = this.$on.outerWidth();
                this._labelWidth = this.$label.outerWidth();
                this.$container.width((this._handleWidth * 2) + this._labelWidth);
                return this.$wrapper.width(this._handleWidth + this._labelWidth);
            };

            BootstrapSwitch.prototype._containerPosition = function(state, callback) {
                if (state == null) {
                    state = this.options.state;
                }
                this.$container.css("margin-left", (function(_this) {
                    return function() {
                        var values;
                        values = [0, "-" + _this._handleWidth + "px"];
                        if (_this.options.indeterminate) {
                            return "-" + (_this._handleWidth / 2) + "px";
                        }
                        if (state) {
                            if (_this.options.inverse) {
                                return values[1];
                            } else {
                                return values[0];
                            }
                        } else {
                            if (_this.options.inverse) {
                                return values[0];
                            } else {
                                return values[1];
                            }
                        }
                    };
                })(this));
                if (!callback) {
                    return;
                }
                return setTimeout(function() {
                    return callback();
                }, 50);
            };

            BootstrapSwitch.prototype._init = function() {
                var init, initInterval;
                init = (function(_this) {
                    return function() {
                        _this.setPrevOptions();
                        _this._width();
                        return _this._containerPosition(null, function() {
                            if (_this.options.animate) {
                                return _this.$wrapper.addClass(_this.options.baseClass + "-animate");
                            }
                        });
                    };
                })(this);
                if (this.$wrapper.is(":visible")) {
                    return init();
                }
                return initInterval = window.setInterval((function(_this) {
                    return function() {
                        if (_this.$wrapper.is(":visible")) {
                            init();
                            return window.clearInterval(initInterval);
                        }
                    };
                })(this), 50);
            };

            BootstrapSwitch.prototype._elementHandlers = function() {
                return this.$element.on({
                    "setPreviousOptions.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            return _this.setPrevOptions();
                        };
                    })(this),
                    "previousState.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            _this.options = _this.prevOptions;
                            if (_this.options.indeterminate) {
                                _this.$wrapper.addClass(_this.options.baseClass + "-indeterminate");
                            }
                            return _this.$element.prop("checked", _this.options.state).trigger("change.bootstrapSwitch", true);
                        };
                    })(this),
                    "change.bootstrapSwitch": (function(_this) {
                        return function(e, skip) {
                            var state;
                            e.preventDefault();
                            e.stopImmediatePropagation();
                            state = _this.$element.is(":checked");
                            _this._containerPosition(state);
                            if (state === _this.options.state) {
                                return;
                            }
                            _this.options.state = state;
                            _this.$wrapper.toggleClass(_this.options.baseClass + "-off").toggleClass(_this.options.baseClass + "-on");
                            if (!skip) {
                                if (_this.$element.is(":radio")) {
                                    $("[name='" + (_this.$element.attr('name')) + "']").not(_this.$element).prop("checked", false).trigger("change.bootstrapSwitch", true);
                                }
                                return _this.$element.trigger("switchChange.bootstrapSwitch", [state]);
                            }
                        };
                    })(this),
                    "focus.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            e.preventDefault();
                            return _this.$wrapper.addClass(_this.options.baseClass + "-focused");
                        };
                    })(this),
                    "blur.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            e.preventDefault();
                            return _this.$wrapper.removeClass(_this.options.baseClass + "-focused");
                        };
                    })(this),
                    "keydown.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            if (!e.which || _this.options.disabled || _this.options.readonly) {
                                return;
                            }
                            switch (e.which) {
                                case 37:
                                    e.preventDefault();
                                    e.stopImmediatePropagation();
                                    return _this.state(false);
                                case 39:
                                    e.preventDefault();
                                    e.stopImmediatePropagation();
                                    return _this.state(true);
                            }
                        };
                    })(this)
                });
            };

            BootstrapSwitch.prototype._handleHandlers = function() {
                this.$on.on("click.bootstrapSwitch", (function(_this) {
                    return function(event) {
                        event.preventDefault();
                        event.stopPropagation();
                        _this.state(false);
                        return _this.$element.trigger("focus.bootstrapSwitch");
                    };
                })(this));
                return this.$off.on("click.bootstrapSwitch", (function(_this) {
                    return function(event) {
                        event.preventDefault();
                        event.stopPropagation();
                        _this.state(true);
                        return _this.$element.trigger("focus.bootstrapSwitch");
                    };
                })(this));
            };

            BootstrapSwitch.prototype._labelHandlers = function() {
                return this.$label.on({
                    "click": function(e) {
                        return e.stopPropagation();
                    },
                    "mousedown.bootstrapSwitch touchstart.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            if (_this._dragStart || _this.options.disabled || _this.options.readonly) {
                                return;
                            }
                            e.preventDefault();
                            e.stopPropagation();
                            _this._dragStart = (e.pageX || e.originalEvent.touches[0].pageX) - parseInt(_this.$container.css("margin-left"), 10);
                            if (_this.options.animate) {
                                _this.$wrapper.removeClass(_this.options.baseClass + "-animate");
                            }
                            return _this.$element.trigger("focus.bootstrapSwitch");
                        };
                    })(this),
                    "mousemove.bootstrapSwitch touchmove.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            var difference;
                            if (_this._dragStart == null) {
                                return;
                            }
                            e.preventDefault();
                            difference = (e.pageX || e.originalEvent.touches[0].pageX) - _this._dragStart;
                            if (difference < -_this._handleWidth || difference > 0) {
                                return;
                            }
                            _this._dragEnd = difference;
                            return _this.$container.css("margin-left", _this._dragEnd + "px");
                        };
                    })(this),
                    "mouseup.bootstrapSwitch touchend.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            var state;
                            if (!_this._dragStart) {
                                return;
                            }
                            e.preventDefault();
                            if (_this.options.animate) {
                                _this.$wrapper.addClass(_this.options.baseClass + "-animate");
                            }
                            if (_this._dragEnd) {
                                state = _this._dragEnd > -(_this._handleWidth / 2);
                                _this._dragEnd = false;
                                _this.state(_this.options.inverse ? !state : state);
                            } else {
                                _this.state(!_this.options.state);
                            }
                            return _this._dragStart = false;
                        };
                    })(this),
                    "mouseleave.bootstrapSwitch": (function(_this) {
                        return function(e) {
                            return _this.$label.trigger("mouseup.bootstrapSwitch");
                        };
                    })(this)
                });
            };

            BootstrapSwitch.prototype._externalLabelHandler = function() {
                var $externalLabel;
                $externalLabel = this.$element.closest("label");
                return $externalLabel.on("click", (function(_this) {
                    return function(event) {
                        event.preventDefault();
                        event.stopImmediatePropagation();
                        if (event.target === $externalLabel[0]) {
                            return _this.toggleState();
                        }
                    };
                })(this));
            };

            BootstrapSwitch.prototype._formHandler = function() {
                var $form;
                $form = this.$element.closest("form");
                if ($form.data("bootstrap-switch")) {
                    return;
                }
                return $form.on("reset.bootstrapSwitch", function() {
                    return window.setTimeout(function() {
                        return $form.find("input").filter(function() {
                            return $(this).data("bootstrap-switch");
                        }).each(function() {
                            return $(this).bootstrapSwitch("state", this.checked);
                        });
                    }, 1);
                }).data("bootstrap-switch", true);
            };

            BootstrapSwitch.prototype._getClasses = function(classes) {
                var c, cls, i, len;
                if (!$.isArray(classes)) {
                    return [this.options.baseClass + "-" + classes];
                }
                cls = [];
                for (i = 0, len = classes.length; i < len; i++) {
                    c = classes[i];
                    cls.push(this.options.baseClass + "-" + c);
                }
                return cls;
            };

            return BootstrapSwitch;

        })();
        $.fn.bootstrapSwitch = function() {
            var args, option, ret;
            option = arguments[0], args = 2 <= arguments.length ? slice.call(arguments, 1) : [];
            ret = this;
            this.each(function() {
                var $this, data;
                $this = $(this);
                data = $this.data("bootstrap-switch");
                if (!data) {
                    $this.data("bootstrap-switch", data = new BootstrapSwitch(this, option));
                }
                if (typeof option === "string") {
                    return ret = data[option].apply(data, args);
                }
            });
            return ret;
        };
        $.fn.bootstrapSwitch.Constructor = BootstrapSwitch;
        return $.fn.bootstrapSwitch.defaults = {
            state: true,
            size: null,
            animate: true,
            disabled: false,
            readonly: false,
            indeterminate: false,
            inverse: false,
            radioAllOff: false,
            onColor: "primary",
            offColor: "default",
            onText: "ON",
            offText: "OFF",
            labelText: "&nbsp;",
            handleWidth: "auto",
            labelWidth: "auto",
            baseClass: "bootstrap-switch",
            wrapperClass: "wrapper",
            onInit: function() {},
            onSwitchChange: function() {}
        };
    })(window.jQuery, window);

}).call(this);