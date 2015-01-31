// Copyright 2014, 2015 Zac Bergquist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SimpleArtist contains basic info about an artist.
type SimpleArtist struct {
	// The name of the artist.
	Name string `json:"name"`
	// The Spotify ID for the artist.
	ID ID `json:"id"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
	// A link to the Web API enpoint providing
	// full details of the artist.
	Endpoint string `json:"href"`
	// Known external URLs for this artist.
	ExternalURLs ExternalURL `json:"external_urls"`
}

// FullArtist provides extra artist data in addition
// to what is provided by SimpleArtist.
type FullArtist struct {
	SimpleArtist
	// The popularity of the artist.  The value will be
	// between 0 and 100, with 100 being the most popular.
	// The artist's popularity is calculated from the
	// popularity of all of the artist's tracks.
	Popularity int `json:"popularity"`
	// A list of genres the artist is associated with.
	// For example, "Prog Rock" or "Post-Grunge".  If
	// not yet classified, the slice is empty.
	Genres []string `json:"genres"`
	// Information about followers of the artist.
	Followers Followers
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

// FindArtist is a wrapper around DefaultClient.FindArtist.
func FindArtist(id ID) (*FullArtist, error) {
	return DefaultClient.FindArtist(id)
}

// FindArtist gets Spotify catalog information for a single
// artist, given that artist's Spotify ID.
func (c *Client) FindArtist(id ID) (*FullArtist, error) {
	uri := baseAddress + "artists/" + string(id)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a FullArtist
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// FindArtists is a wrapper around DefaultClient.FindArtists.
func FindArtists(ids ...ID) ([]*FullArtist, error) {
	return DefaultClient.FindArtists(ids...)
}

// FindArtists gets spotify catalog information for several
// artists based on their Spotify IDs.  It supports up to
// 50 artists in a single call.  Artists are returned in the
// order requested.  If an artist is not found, that position
// in the result will be nil.  Duplicate IDs will result in
// duplicate artists in the result.
func (c *Client) FindArtists(ids ...ID) ([]*FullArtist, error) {
	uri := baseAddress + "artists?ids=" + strings.Join(toStringSlice(ids), ",")
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a struct {
		Artists []*FullArtist
	}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a.Artists, nil
}

// ArtistsTopTracks is a wrapper around DefaultClient.ArtistTopTracks.
func ArtistsTopTracks(artistID ID, country string) ([]FullTrack, error) {
	return DefaultClient.ArtistsTopTracks(artistID, country)
}

// ArtistsTopTracks gets Spotify catalog information about
// an artist's top tracks in a particular country.  It returns
// a maximum of 10 tracks.  The country is specified as an
// ISO 3166-1 alpha-2 country code.
func (c *Client) ArtistsTopTracks(artistID ID, country string) ([]FullTrack, error) {
	uri := baseAddress + "artists/" + string(artistID) + "/top-tracks?country=" + country
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var t struct {
		Tracks []FullTrack `json:"tracks"`
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	return t.Tracks, nil
}

// FindRelatedArtists is a wrapper around DefaultClient.FindRelatedArtists.
func FindRelatedArtists(id ID) ([]FullArtist, error) {
	return DefaultClient.FindRelatedArtists(id)
}

// FindRelatedArtists gets Spotify catalog information about
// artists similar to a given artist.  Similarity is based on
// analysis of the Spotify community's listening history.
// This function returns up to 20 artists that are considered
// related to the specified artist.
func (c *Client) FindRelatedArtists(id ID) ([]FullArtist, error) {
	uri := baseAddress + "artists/" + string(id) + "/related-artists"
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a struct {
		Artists []FullArtist `json:"artists"`
	}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a.Artists, nil
}

// ArtistAlbums is a wrapper around DefaultClient.ArtistAlbums.
func ArtistAlbums(artistID ID) (*SimpleAlbumPage, error) {
	return DefaultClient.ArtistAlbums(artistID)
}

// ArtistAlbums gets Spotify catalog information about an artist's albums.
// It is equivalent to ArtistAlbumsOpt(artistID, nil).
func (c *Client) ArtistAlbums(artistID ID) (*SimpleAlbumPage, error) {
	return c.ArtistAlbumsOpt(artistID, nil, nil)
}

// ArtistAlbumsOpt is a wrapper around DefaultClient.ArtistAlbumsOpt
func ArtistAlbumsOpt(artistID ID, options *Options, t *AlbumType) (*SimpleAlbumPage, error) {
	return DefaultClient.ArtistAlbumsOpt(artistID, options, t)
}

// ArtistAlbumsOpt is just like ArtistAlbums, but it accepts optional parameters
// to filter and sort the result.
//
// The AlbumType argument can be used to find a particular type of album.  Search
// for multiple types by OR-ing the types together.
func (c *Client) ArtistAlbumsOpt(artistID ID, options *Options, t *AlbumType) (*SimpleAlbumPage, error) {
	uri := baseAddress + "artists/" + string(artistID) + "/albums"
	// add optional query string if options were specified
	if options != nil {
		values := url.Values{}
		if t != nil {
			values.Set("album_type", t.encode())
		}
		if options.Country != nil {
			values.Set("market", *options.Country)
		} else {
			// if the market is not specified, Spotify will likely return a lot
			// of duplicates (one for each market in which the album is available)
			// - prevent this behavior by falling back to the US by default
			// TODO: would this ever be the desired behavior?
			values.Set("market", CountryUSA)
		}
		if options.Limit != nil {
			values.Set("limit", strconv.Itoa(*options.Limit))
		}
		if options.Offset != nil {
			values.Set("offset", strconv.Itoa(*options.Offset))
		}
		if query := values.Encode(); query != "" {
			uri += "?" + query
		}
	}
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var p rawPage
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	var result SimpleAlbumPage
	err = json.Unmarshal([]byte(p.Items), &result.Albums)
	if err != nil {
		return nil, err
	}
	result.Endpoint = p.Endpoint
	result.Limit = p.Limit
	result.Offset = p.Offset
	result.Total = p.Total
	result.Previous = p.Previous
	result.Next = p.Next
	return &result, nil
}
