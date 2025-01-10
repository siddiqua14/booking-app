<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Property Listings</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css">
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

    <div id="propertyGrid" class="grid">
        {{range .Properties}}
        <div class="property-card">
            <div class="property-image">
                <img src="{{if .HeroImage}}{{.HeroImage}}{{else}}/static/images/placeholder.jpg{{end}}" alt="{{.HotelName}}">
            </div>
            <div class="property-details">
                <div class="rating">
                    <span class="rating-value">{{printf "%.1f" .Rating}}</span>
                    <span class="review-count">({{.ReviewCount}} Reviews)</span>
                </div>
                <h3 class="property-title">
                    <a href="/property-details?id={{.HotelID}}" target="_blank">{{.HotelName}}</a>
                </h3>
                <div class="location">{{.Location}}</div>
                <div class="price">From {{.Price}}</div>
                <button class="availability-button" onclick="window.open('/property-details?id={{.HotelID}}', '_blank')">View Availability</button>
            </div>
        </div>
        {{end}}
    </div>

    <script src="/static/js/property.js"></script>
</body>
</html>