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
<div class="property-details">
    <div class="property-header">
        <h1>{{.Property.HotelName}}</h1>
        <div class="property-meta">
            <span class="rating">★ {{.Property.Rating}} ({{.Property.ReviewCount}} Reviews)</span>
            <span class="details">{{.Property.Bedroom}} Bedroom</span>
            <span class="details">{{.Property.NumBaths}} Bathroom</span>
            <span class="details">{{.Property.Guests}} Guests</span>
        </div>
    </div>

    <div class="property-images">
        <div class="main-image">
            <img src="{{.Details.ImageUrl1}}" alt="Property Main Image">
        </div>
        <div class="image-grid">
            <img src="{{.Details.ImageUrl2}}" alt="Property Image">
            <img src="{{.Details.ImageUrl3}}" alt="Property Image">
            <img src="{{.Details.ImageUrl4}}" alt="Property Image">
            <img src="{{.Details.ImageUrl5}}" alt="Property Image">
        </div>
    </div>

    <div class="property-description">
        <h2>About this property</h2>
        <p>{{.Details.Description}}</p>
    </div>
</div>
</body>
</html>