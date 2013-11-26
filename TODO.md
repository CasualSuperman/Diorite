TODO
====

* Figure out what kind of GUI we're using.
	* QML
	* GTK
	* QT
	* Web with embedded browser
* ~~Finish set legality filters.~~
* Finish other filters.
* ~~Get a ban list/restricted list in a downloadable format.~~
    * ~~searchthecity.me can do this! We should probably set up a service to ping their data occassionally and host it locally, since they don't send modified headers.~~
	* ~~Need to find a way to discover differences without saving lists. Use a hash of the payload, maybe?~~
* ~~Should rehost mtgjson.com's source before release.~~ (But give proper credit)
* Find a card price API?
* ~~Get a different db storage location for Windows (OSX too?) since it doesn't support the homedir method.~~
    * Confirm the storage locations for OSX and Windows are correct.
* Rewrite the gob info to not require custom types, just use GobEncoders and GobDecoders for special types.
    * Blockedon [Issue 6737](https://code.google.com/p/go/issues/detail?id=6737).
* ~~Unify all the []Card and []*Cards hanging out in various places.~~
* Check our pointer usage in various locations.
    * Diorite/Multiverse.Sets
	* Diorite/Multiverse.Cardlist.Add
	* Others? Go hunting!
* ~~Rewrite server multiverse downloading to be more parallel and detect differences in banlists.~~
* ~~Add supertypes back in to cards.~~
* ~~Create enums for supertypes and types.~~
* ~~Add a pretty print method for cards.~~
