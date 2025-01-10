<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Property Details</title>
    <link rel="stylesheet" href="/static/css/property.css">
</head>
<body>
    <header class="header-container">
        <div class="header">
            <a href="/v1/property/list">
                <img src="https://static.rentbyowner.com/release/28.0.6/static/images/sites/rentbyowner.com/header_logo.svg" alt="rentbyowner logo">
            </a>
            <div class="search-container">
                <input type="text" id="location" placeholder="Dubai, Dubai, United Arab Emirates..." />
                <input type="date" id="date" />
                <select id="guests">
                    <option value="">Select guests</option>
                    <option value="1">1 Guest</option>
                    <option value="2">2 Guests</option>
                    <option value="3">3 Guests</option>
                    <option value="4">4 Guests</option>
                    <option value="5">5 Guests</option>
                </select>
            </div>
        </div>
    </header>

    <main>
        <div class="property-detail">
            <h1>{{.Property.HotelName}}</h1>
            <p><strong>Location:</strong> {{.Property.Location}}</p>
            <p><strong>Price:</strong> {{.Property.Price}}</p>
            <p><strong>Type:</strong> {{.Property.PropertyType}}</p>
            <p><strong>Guests:</strong> {{.Property.Guests}}</p>
            <p><strong>Rating:</strong> {{.Property.Rating}} ({{.Property.ReviewCount}} Reviews)</p>
            <p><strong>Beds:</strong> {{.Property.NumBeds}}</p>
            <p><strong>Bedrooms:</strong> {{.Property.NumBedR}}</p>
            <p><strong>Bathrooms:</strong> {{.Property.NumBaths}}</p>
            <p><strong>Bedrooms:</strong> {{.Property.Bedroom}}</p>
            <p><strong>Amenities:</strong> {{range .Property.Amenities}}{{.}}, {{end}}</p>
            <div class="images">
                {{range .Property.Images}}
                <img src="{{.}}" alt="Property Image" class="property-image">
                {{end}}
            </div>
            <p>{{.Property.Description}}</p>
        </div>
    </main>
    <script src="/static/js/property.js"></script>
</body>
</html>