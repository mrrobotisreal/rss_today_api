# Google News RSS Processing Improvements

## Overview

This document explains the improvements made to handle Google News RSS feeds and address the issues you were experiencing with encoded URLs and HTML content.

## What Was the Problem?

### Before the Improvements:
1. **Encoded URLs**: Google News RSS feeds return encoded redirect URLs like `https://news.google.com/rss/articles/CBMi...` instead of direct article URLs
2. **HTML Content**: Titles and descriptions contained HTML entities and tags that weren't being properly cleaned
3. **Source Names in Titles**: Google News titles often included source names at the end (e.g., "Article Title - BBC News")
4. **Tracking Parameters**: URLs contained unnecessary tracking parameters

### What You Were Seeing:
- HTML list format output instead of clean article data
- Long encoded Google News URLs instead of direct article links
- HTML entities like `&amp;`, `&quot;` in titles and descriptions

## What Is "Autodiscovery"?

**Autodiscovery** refers to automatically finding RSS feed URLs from a website's homepage by looking for `<link>` tags in the HTML head. This is **NOT** what you need for Google News.

Your issue was specifically with **URL decoding** and **content cleaning**, not autodiscovery.

## Improvements Made

### 1. Enhanced RSS Feed Processing (`internal/services/rss_feed.go`)

#### New Content Cleaning:
- **HTML Entity Decoding**: Converts `&amp;` → `&`, `&quot;` → `"`, etc.
- **HTML Tag Removal**: Strips `<a>`, `<font>`, and other HTML tags
- **Source Name Removal**: Removes trailing source names from titles
- **Whitespace Cleanup**: Normalizes spacing

#### Better URL Processing:
- **Google News URL Detection**: Identifies encoded Google News URLs
- **URL Decoding**: Attempts to decode Google News URLs to original article URLs
- **Tracking Parameter Removal**: Removes UTM and other tracking parameters

### 2. Google News URL Decoder (`internal/services/google_news_decoder.go`)

A specialized service for decoding Google News encoded URLs:

#### Features:
- **Base64 Decoding**: Attempts to decode the encoded part of Google News URLs
- **HTTP Redirect Following**: Follows redirects to find original URLs
- **Batch Processing**: Can decode multiple URLs efficiently
- **Error Handling**: Gracefully falls back to original URL if decoding fails

#### Methods Used:
1. **Base64 Analysis**: Extracts and decodes the base64-encoded portion of Google News URLs
2. **HTTP Head Requests**: Follows redirects to discover the target URL
3. **Pattern Matching**: Uses regex to find URL patterns in decoded data

### 3. Test Endpoint (`/test/google-news`)

A new endpoint to help you see the improvements in action:

```bash
GET http://localhost:3001/test/google-news
```

#### Response Format:
```json
{
  "source": "Google News",
  "total_articles": 100,
  "processed_articles": [
    {
      "original_title": "Article Title - BBC News&amp;nbsp;",
      "cleaned_title": "Article Title",
      "original_description": "<p>Description with &quot;quotes&quot;</p>",
      "cleaned_description": "Description with \"quotes\"",
      "original_link": "https://news.google.com/rss/articles/CBMi...",
      "processed_link": "https://www.bbc.com/news/article-123",
      "is_google_news_url": true,
      "url_decoded": true
    }
  ]
}
```

## How to Use

### 1. Test the Improvements

Start your API and call the test endpoint:

```bash
# Start your API
go run cmd/api/main.go

# Test the improvements
curl http://localhost:3001/test/google-news
```

### 2. Monitor Regular RSS Processing

The improvements are automatically applied to all RSS processing. You can trigger manual monitoring:

```bash
curl -X POST http://localhost:3001/api/monitor/trigger \
  -H "Authorization: Bearer YOUR_FIREBASE_TOKEN"
```

### 3. View Cleaned Articles

Get articles through the API to see the cleaned content:

```bash
curl http://localhost:3001/api/articles \
  -H "Authorization: Bearer YOUR_FIREBASE_TOKEN"
```

## Technical Implementation Details

### URL Decoding Process

1. **Detection**: Check if URL contains `news.google.com` and `/articles/`
2. **Base64 Extraction**: Extract the encoded part after `/articles/`
3. **Decoding**: Attempt base64 decoding (both URL-safe and standard)
4. **Pattern Matching**: Look for HTTP(S) URLs in decoded data
5. **Validation**: Ensure decoded URL is not another Google URL
6. **Fallback**: If decoding fails, keep original URL

### Content Cleaning Process

1. **HTML Entity Decoding**: Use Go's `html.UnescapeString()`
2. **HTML Tag Removal**: Regex-based tag stripping
3. **Source Name Removal**: Remove trailing ` - Source Name` patterns
4. **Whitespace Normalization**: Clean up extra spaces and newlines

### Configuration

The Google News decoder is automatically initialized for sources that:
- Have "Google News" in the name, OR
- Have `news.google.com` in the RSS URL

## Expected Results

### Before:
```html
<li><a href="https://news.google.com/rss/articles/CBMi...">
Article Title - BBC News</a>&nbsp;&nbsp;<font color="#6f6f6f">BBC News</font></li>
```

### After:
```json
{
  "title": "Article Title",
  "description": "Clean article description",
  "link": "https://www.bbc.com/news/actual-article-url",
  "source": "Google News"
}
```

## Limitations and Considerations

### URL Decoding Limitations:
1. **Rate Limiting**: Google may rate-limit decoding requests
2. **New Encoding**: Google occasionally changes their encoding scheme
3. **Success Rate**: Not all URLs may be decodable due to various factors

### Recommendations:
1. **Monitor Logs**: Check for "Successfully decoded" vs "Could not decode" messages
2. **Rate Limiting**: The decoder includes delays to be respectful to Google's servers
3. **Fallback Handling**: Original URLs are preserved if decoding fails

## Troubleshooting

### If URLs Are Not Being Decoded:
1. Check logs for "Google News decoder initialized" messages
2. Verify the source name contains "Google News" or URL contains `news.google.com`
3. Check for rate limiting or network issues

### If Content Is Still Not Clean:
1. Verify HTML entities are being decoded
2. Check for new HTML patterns that need handling
3. Review the `CleanContent` function for additional cleaning rules

## Future Enhancements

### Potential Improvements:
1. **Advanced Decoding**: Implement more sophisticated Google News URL decoding methods
2. **Caching**: Cache decoded URLs to reduce API calls
3. **Source Detection**: Better source extraction from decoded URLs
4. **Custom Cleaning**: Source-specific content cleaning rules

### Monitoring:
1. **Success Metrics**: Track URL decoding success rates
2. **Performance**: Monitor decoding response times
3. **Quality**: Track content cleaning effectiveness