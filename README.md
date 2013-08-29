Ferret
======
##An optimized substring search engine written in Go.
Ferret makes use of a combination of an Inverted Index and a Suffix Array to allow log-time lookups with a relatively small memory footprint.
Also incorporates error-correction (Levenshtein distance 1) and simple Unicode-to-ASCII conversion.
Allows for arbitrary sorting functions
Allows you to map arbitrary data to your results, and quickly update this data.

***Author:*** Mark Canning 
***Developed at/for:*** Tamber - http://www.tamber.com/

About Tamber
------------
Tamber also has this really cool recommendation engine for music (also development by me) that prioritizes up-and-coming artists, so that it doesn't succomb to popularity biases, but still produces great recommendations! Make sure to check us out at http://www.tamber.com or https://itunes.apple.com/us/app/tamber-concerts/id658240483

Installing
----------
To Install: go get github.com/argusdusty/Ferret
To Update: go get -u github.com/argusdusty/Ferret
To Use: import "github.com/argusdusty/Ferret"

Performance
-----------
Uses linear memory (~10-18 bytes per character)
Searches performed in log time with the number of characters in the dictionary.
Sorted searches can be slow, taking ~linear time with the number of matches, rather than linear time with the results limit.
Initialization takes linearithmic (ln(n)*n) time (being a sorting algorithm)

The code is meant to be as fast as possible for a substring dictionary search, and as such is best suited for medium-large dictionaries with ~1-100 million total characters. I've timed 10s initialization for 3.5 million characters on a modern CPU, and 10us search time (4000us with error-correction), so this system is capable of ~100,000 queries per second on a single processor.


Sample usage
------------

###Initializing the search engine:
```go
// Allows for exact (case-sensitive) substring searches over a list of songs mapping their respective artists, allowing sorting by the song popularity
SearchEngine := ferret.New(Songs, Artists, SongPopularities, func(s string) []byte { return []byte(s) })

// Allows for lowercase-ASCII substring searches over a list of songs mapping their respective artists, allowing sorting by the song popularity
SearchEngine := ferret.New(Songs, Artists, SongPopularities, ferret.UnicodeToLowerASCII)

// Allows for lowercase-ASCII substring searches over a list of artists, allowing sorting by the artist popularity
SearchEngine := ferret.New(Artists, Artists, ArtistPopularities, ferret.UnicodeToLowerASCII)
```
		
###Inserting a new element into the search engine:
```go
// Add a song to an existing SearchEngine, written by Artist, and with popularity SongPopularity
SearchEngine.Insert(Song, Artist, SongPopularity)
```

###Performing simple unsorted substring search:
```go
// For songs - returns a list of up to 25 artists of the matching songs, and the song popularities
SearchEngine.Query(SongQuery, 25)
```
	
###Performing a sorted substring search:
```go
// For songs - returns a list of up to 25 artists of the matching songs, and the song popularities, sorted by the song popularities
// assuming the song popularities are float64s
SearchEngine.SortedQuery(SongQuery, 25, func(s string, v interface{}, l int, i int) float64 { return v.(float64) })
```

###More examples	
Check out example.go and dictionaryexample.go for more example usage.
