	// Globals variables     
      
    var map;
    
        // 	An array containing objects with information about the locations of bisness biz.
        var biz = [];
    
          // Create a new blank array for all the listing markers.
          var markers = [];
          var markersgroup = [];
    
          function initMap() {
            // Create a styles array to use with the map.
            var styles = [
              {
                featureType: 'water',
                stylers: [
                  { color: '#19a0d8' }
                ]
              },{
                featureType: 'administrative',
                elementType: 'labels.text.stroke',
                stylers: [
                  { color: '#ffffff' },
                  { weight: 6 }
                ]
              },{
                featureType: 'administrative',
                elementType: 'labels.text.fill',
                stylers: [
                  { color: '#e85113' }
                ]
              },{
                featureType: 'road.highway',
                elementType: 'geometry.stroke',
                stylers: [
                  { color: '#efe9e4' },
                  { lightness: -40 }
                ]
              },{
                featureType: 'transit.station',
                stylers: [
                  { weight: 9 },
                  { hue: '#e85113' }
                ]
              },{
                featureType: 'road.highway',
                elementType: 'labels.icon',
                stylers: [
                  { visibility: 'off' }
                ]
              },{
                featureType: 'water',
                elementType: 'labels.text.stroke',
                stylers: [
                  { lightness: 100 }
                ]
              },{
                featureType: 'water',
                elementType: 'labels.text.fill',
                stylers: [
                  { lightness: -100 }
                ]
              },{
                featureType: 'poi',
                elementType: 'geometry',
                stylers: [
                  { visibility: 'on' },
                  { color: '#f0e4d3' }
                ]
              },{
                featureType: 'road.highway',
                elementType: 'geometry.fill',
                stylers: [
                  { color: '#efe9e4' },
                  { lightness: -25 }
                ]
              }
            ];
    
            // Constructor creates a new map - only center and zoom are required.
            map = new google.maps.Map(document.getElementById('map'), {
              center: {lat: 6.463248, lng: 3.386347},
              zoom: 13,
              styles: styles,
              mapTypeControl: false
            });
    
            // These are the real estate listings that will be shown to the user.
             // Normally we'd have these in a database instead.
            biz = [
              idumota = [ 
                clothing = [
                    {title: 'Park Ave Penthouse', location: {lat: 6.463248, lng: 3.386347}},
                    {title: 'Chelsea Loft', location: {lat: 6.463348, lng: 3.386347}},
                    {title: 'Union Square Open Floor Plan', location: {lat: 6.463248, lng: 3.386547}},
                    {title: 'East Village Hip Studio', location: {lat: 6.463258, lng: 3.386447}},
                    {title: 'TriBeCa Artsy Bachelor Pad', location: {lat: 6.463141, lng: 3.386340}},
                    {title: 'Chinatown Homey Space', location: {lat: 6.463348, lng: 3.386447}}
                ], 
                jewelries = [
                    {title: 'Park Ave Penthouse', location: {lat: 6.463248, lng: 3.386447}},
                    {title: 'Chelsea Loft', location: {lat: 6.463348, lng: 3.386247}},
                    {title: 'Union Square Open Floor Plan', location: {lat: 6.463548, lng: 3.386547}},
                    {title: 'East Village Hip Studio', location: {lat: 6.463158, lng: 3.386347}},
                    {title: 'TriBeCa Artsy Bachelor Pad', location: {lat: 6.463141, lng: 3.386390}},
                    {title: 'Chinatown Homey Space', location: {lat: 6.463548, lng: 3.386447}}
                ] 
              ],
              hotel = [ 
                lagos = [
                  {title: 'Park Ave Penthouse', location: {lat: 6.463248, lng: 3.386447}},
                  {title: 'Chelsea Loft', location: {lat: 6.463348, lng: 3.386247}},
                  {title: 'Union Square Open Floor Plan', location: {lat: 6.463548, lng: 3.386547}},
                  {title: 'East Village Hip Studio', location: {lat: 6.463158, lng: 3.386347}},
                  {title: 'TriBeCa Artsy Bachelor Pad', location: {lat: 6.463141, lng: 3.386390}},
                  {title: 'Chinatown Homey Space', location: {lat: 6.463548, lng: 3.386447}}
                ],
                imo = [
    
                ],
              ],
            ];
    
            var largeInfowindow = new google.maps.InfoWindow();
           
    
            // Style the markers a bit. This will be our listing marker icon.
            var defaultIcon = makeMarkerIcon('0091ff');
    
            // Create a "highlighted location" marker color for when the user
            // mouses over the marker.
            var highlightedIcon = makeMarkerIcon('FFFF24');
    
            var count = 0
            // The following group uses the location array to create an array of markers on initialize.
            for (var val in biz) {
            for (var val1 in biz[val]) {
            for (var val2 in biz[val][val1]) {
              // Get the position from the location array.
              var position = biz[val][val1][val2].location;
              var title = biz[val][val1][val2].title;
              // Create a marker per location, and put into markers array.
              var marker = new google.maps.Marker({
                position: position,
                title: title,
                animation: google.maps.Animation.DROP,
                icon: defaultIcon,
                id: count
              });
              // Push the marker to our array of markers.
              markers.push(marker);
              // Create an onclick event to open the large infowindow at each marker.
              marker.addListener('click', function() {
                populateInfoWindow(this, largeInfowindow);
              });
         // Two event listeners - one for mouseover, one for mouseout,
              // to change the colors back and forth.
              marker.addListener('mouseover', function() {
                this.setIcon(highlightedIcon);
              });
              marker.addListener('mouseout', function() {
                this.setIcon(defaultIcon);
              }); 
              count++    
            }
          markersgroup[val1] = markers;
          markers = [];
        }}
       
            document.getElementById('show-listings').addEventListener('click', showListings);
            document.getElementById('hide-listings').addEventListener('click', hideListings);
         
               }
    
          // This function populates the infowindow when the marker is clicked. We'll only allow
          // one infowindow which will open at the marker that is clicked, and populate based
          // on that markers position.
          function populateInfoWindow(marker, infowindow) {
            // Check to make sure the infowindow is not already opened on this marker.
            if (infowindow.marker != marker) {
      // Clear the infowindow content to give the streetview time to load.
              infowindow.setContent('');
              infowindow.marker = marker;
              // Make sure the marker property is cleared if the infowindow is closed.
              infowindow.addListener('closeclick', function() {
               infowindow.marker = null;
              });
              var streetViewService = new google.maps.StreetViewService();
              var radius = 50;
              // In case the status is OK, which means the pano was found, compute the
              // position of the streetview image, then calculate the heading, then get a
              // panorama from that and set the options
              function getStreetView(data, status) {
                if (status == google.maps.StreetViewStatus.OK) {
                  var nearStreetViewLocation = data.location.latLng;
                  var heading = google.maps.geometry.spherical.computeHeading(
                    nearStreetViewLocation, marker.position);
                    infowindow.setContent('<div>' + marker.title + '</div><div id="pano"></div>');
                    var panoramaOptions = {
                      position: nearStreetViewLocation,
                      pov: {
                        heading: heading,
                        pitch: 30
                      }
                    };
                  var panorama = new google.maps.StreetViewPanorama(
                    document.getElementById('pano'), panoramaOptions);
                } else {
                  infowindow.setContent('<div>' + marker.title + '</div>' +
                    '<div>No Street View Found</div>');
                }
              }
              // Use streetview service to get the closest streetview image within
              // 50 meters of the markers position
              streetViewService.getPanoramaByLocation(marker.position, radius, getStreetView);
              // Open the infowindow on the correct marker.
              infowindow.open(map, marker);
            }
          }
    
      // This function will loop through the markers array and display them all.
          function showListings() {
            var bounds = new google.maps.LatLngBounds();
            // Extend the boundaries of the map for each marker and display the marker
            //for (var i = 0; i < markers.length; i++) {
               var val = 0;
                for (var val1 in markersgroup[val]) {
              markersgroup[val][val1].setMap(map);
              bounds.extend( markersgroup[val][val1].position);
          }
            map.fitBounds(bounds);
          }
    
          // This function will loop through the listings and hide them all.
          function hideListings() {
             var val = 0;
            for (var val1 in markersgroup[val]) {
             markersgroup[val][val1].setMap(null);
            }
          }
    
          // This function takes in a COLOR, and then creates a new marker
          // icon of that color. The icon will be 21 px wide by 34 high, have an origin
          // of 0, 0 and be anchored at 10, 34).
          function makeMarkerIcon(markerColor) {
            var markerImage = new google.maps.MarkerImage(
              'http://chart.googleapis.com/chart?chst=d_map_spin&chld=1.15|0|'+ markerColor +
              '|40|_|%E2%80%A2',
              new google.maps.Size(21, 34),
              new google.maps.Point(0, 0),
              new google.maps.Point(10, 34),
              new google.maps.Size(21,34));
            return markerImage;
          }
    function openNav() {
         document.getElementById("mySlide").style.width = "250px";
    }
    
    function closeNav() {
        document.getElementById("mySlide").style.width = "0";
    }
    
    function sopenNav() {
         document.getElementById("mySidenav").style.width = "250px";
    }
    
    function scloseNav() {
        document.getElementById("mySidenav").style.width = "6%";
    }
    