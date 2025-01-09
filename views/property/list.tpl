{{define "property_list"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Property Listings</title>
    <link rel="stylesheet" href="/static/css/property.css">
</head>
<body>
     <header class="header">
        <div class="search-container">
            <input type="text" class="search-input" id="locationSearch" placeholder="Enter location...">
            <button class="search-button" onclick="searchProperties()">Search</button>
        </div>
    </header>

    <div class="breadcrumb">
        <span id="locationBreadcrumb">Home</span>
    </div>

    <div class="property-grid" id="propertyGrid">
        <!-- Properties will be loaded here via JavaScript -->
    </div>

    <script>
        function searchProperties() {
            const location = document.getElementById('locationSearch').value;
            if (!location) return;

            // Update breadcrumb
            document.getElementById('locationBreadcrumb').innerHTML = 
                `<a href="#">Home</a> > ${location}`;

            // Fetch properties
            fetch(`/v1/property/list?location=${encodeURIComponent(location)}`)
                .then(response => response.json())
                .then(data => displayProperties(data))
                .catch(error => console.error('Error:', error));
        }

        function displayProperties(properties) {
            const grid = document.getElementById('propertyGrid');
            grid.innerHTML = '';

            properties.forEach(property => {
                const card = `
                    <div class="property-card">
                        <img src="/static/images/${property.HotelID}/1.jpg" alt="${property.HotelName}" class="property-image">
                        <div class="property-info">
                            <div class="property-price">From ${property.Price}</div>
                            <div class="property-rating">
                                <span class="rating-badge">${property.Rating.toFixed(1)}</span>
                                <span>(${property.ReviewCount} Reviews)</span>
                            </div>
                            <h3>${property.HotelName}</h3>
                            <p>${property.PropertyType} • ${property.Guests} guests</p>
                            <div class="amenities">${formatAmenities(property.Amenities)}</div>
                            <button class="view-button" onclick="viewProperty('${property.HotelID}')">
                                View Availability
                            </button>
                        </div>
                    </div>
                `;
                grid.innerHTML += card;
            });
        }

        function formatAmenities(amenities) {
            // Assuming amenities is a comma-separated string
            return amenities.split(',').slice(0, 3).join(' • ');
        }

        function viewProperty(hotelId) {
            window.location.href = `/property/details/${hotelId}`;
        }

        // Load initial properties on page load
        document.addEventListener('DOMContentLoaded', () => {
            searchProperties();
        });
    </script>
</body>
</html>