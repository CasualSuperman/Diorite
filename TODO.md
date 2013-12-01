TODO
====

* Web with embedded browser
	* Find/write an embedded browser.
	* Find a good way to embed our web files instead of shipping them in a subdir.
* Finish other filters.
* ~~Should rehost mtgjson.com's source before release.~~ (But give proper credit)
* Find a card price API?
* Confirm the storage locations for OSX and Windows are correct.
	* Can confirm the location for Windows at least *works*.
* Rewrite the gob info to not require custom types, just use GobEncoders and GobDecoders for special types.
    * Blockedon [Issue 6737](https://code.google.com/p/go/issues/detail?id=6737).
* Check our pointer usage in various locations.
    * Diorite/Multiverse.Sets
	* Diorite/Multiverse.Cardlist.Add
	* Others? Go hunting!
* Remember to use -ldflags="-H windowsgui" when building for windows.
* Find a way to embed the web files into our final binary.
	* go-bindata is a way to do this, but isn't automated or anything.
	* We may want a custom tool for this
		* We should probably search for one first though if we decide not to script go-bindata
* The global multiverse variable in the web package feels hackish. I'd like an alternative.
	* Creating a server type that holds it internally should probably work.
* Add a console to the server application so we can force rechecks, etc.
* Add tappedout integration?
* Pull images from Gatherer?
* Find a way to test the web interface.
	* Maybe that go-webkit2 or whatever I saw on reddit.
