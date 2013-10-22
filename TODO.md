TODO
====

* Make card searching handle multiple words.
    * This means making searches for multiple partial words work proprly
	* Also means sorting based on a different closeness metric at the end
	* N-Gram similarity?
* Figure out what kind of GUI we're using.
	* QML
	* GTK
	* QT
	* Web with embedded browser
* ~~Retool multiverse to use less memory?~~
    * ~~Use a [skiplist](https://code.google.com/p/go-wiki/wiki/Projects#Lists) for multiverseId -> Card~~
* Get a ban list/restricted list in a downloadable format.
    * searchthecity.me can do this! We should probably set up a service to ping their data occassionally and host it locally, since they don't send modified headers.
* Should rehost mtgjson.com's source before release. (But give proper credit)
* Find a card price API?
* Find a way to pull the json logic out of the multiverse package
  The json is implementation specific but there's not really a good way to make a multiverse automatically without knowing what format it's coming in as.
  Find an interface to abstract it, even if it's slower than operating on a known structure?
  Is this necessary? Odds are we won't switch data providers and this will only make things slower.
* Get a different db storage location for Windows (OSX too?) since it doesn't support the homedir method.
* ~~Add compression to the gob'd multiverse cache.~~
* ~~Store our Multiverse data instead of the JSON.~~
