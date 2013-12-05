Node.prototype.on = function(type, listener, useCapture) {
	this.addEventListener(type, listener, useCapture);
	return this;
};
Node.prototype.empty = function() {
	var last;
	while (last = this.lastChild) {
		this.removeChild(last);
	}
	return this;
};
Node.prototype.clear = function() {
	this.value = "";
};
Node.prototype.text = function(content) {
	if (content === undefined) {
		return this.textContent;
	}
	this.textContent = content;
	return this;
}
Node.prototype.data = function(key, val) {
	if (arguments.length == 1) {
		return this.dataset[key];
	}
	this.dataset[key] = val;
	return this;
}
Node.prototype.append = function() {
	var n = this;
	Array.prototype.forEach.call(arguments, function(arg) {
		if (Array.isArray(arg) || arg instanceof NodeList) {
			arg.forEach(function(node) {
				n.appendChild(node);
			});
		} else if (arg instanceof Node) {
			n.appendChild(arg);
		} else {
			// wat
		}
	});
	return this;
}
NodeList.prototype.each = function() {
	Array.prototype.forEach.apply(this, arguments);
	return this;
};
NodeList.prototype.on = function() {
	this.each(function(elem) {
		elem.on.apply(elem, arguments);
	});
};
window.$ = function(node, selector) {
	var elems = document.querySelectorAll(node, selector);
	return elems.length == 1 ? elems[0] : elems;
}
window.$.create = function(elem) {
	return document.createElement(elem);
}