TODO
====

* Parse rules text.
* ~~Find a replacement for searchthecity.me now that it's just an info site.~~
    * ~~Port the code?~~
* Web with embedded browser
	* Find/write an embedded browser.
	* Find a good way to embed our web files instead of shipping them in a subdir.
		* ~~go-bindata is a way to do this, but isn't automated or anything.~~
		* We may want a custom tool for this
			* We should probably search for one first though if we decide not to script go-bindata
		* github.com/carbocation/gotogether ?
* Finish other filters.
* ~~Should rehost mtgjson.com's source before release.~~ (But give proper credit)
    * ~~Use their new zip file to save their bandwidth.~~
* Find a card price API?
* Confirm the storage locations for OSX and Windows are correct.
	* Can confirm the location for Windows at least *works*.
* Rewrite the gob info to not require custom types, just use GobEncoders and GobDecoders for special types.
    * ~~Blockedon [Issue 6737](https://code.google.com/p/go/issues/detail?id=6737).~~
* Check our pointer usage in various locations.
    * Diorite/Multiverse.Sets
	* Diorite/Multiverse.Cardlist.Add
	* Others? Go hunting!
	* Made sure all the structs are using pointers properly.
	* Still need to check on arguments/method receivers.
* Remember to use -ldflags="-H windowsgui" when building for windows.
* ~~The global multiverse variable in the web package feels hackish. I'd like an alternative.~~
	* ~~Creating a server type that holds it internally should probably work.~~
* Add a console to the server application so we can force rechecks, etc.
* Add tappedout integration?
* Pull images from Gatherer?
* Find a way to test the web interface.
	* Maybe that go-webkit2 or whatever I saw on reddit.
