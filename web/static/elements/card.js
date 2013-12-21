var MtgCardPrototype = Object.create(HTMLElement.prototype);

MtgCardPrototype.createdCallback = function() {
	var shadow = this.createShadowRoot();
	var copy = document.querySelector("template#mtg_card").cloneNode(true);
	shadow.append(copy.content);
	var id = this.getAttribute("multiverse-id");
	if (id) {
		//shadow.innerHTML = ;
	}
};

var MtgCard = document.register('mtg-card', {
	prototype: MtgCardPrototype
});