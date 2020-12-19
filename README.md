# WebCrawler

WebCrawler is a web application that provides its users with some intel about a given URL.

Test Online: [webcrawler.geekembly.com](http://webcrawler.geekembly.com)

Webcrawler returns information about
* Page title
* HTML version
* Number of links
* Number of internal links
* Number of external links
* Number of total links	
* Number of broken links
* Existance of a Login form in the page

## Project assumptions:

### Redirections
Some links when access are redirected to another link. for example, Yahoo, automatically directs all requsts for [www.yahoo.com](http://www.yahoo.com) to [https://de.yahoo.com/?p=us](https://de.yahoo.com/?p=us) in Germany. WebCrawler follows these redirections and show information about the target link.

### Irrigular links
Irrigular links such as `mailto:a@a.com` or `tel:1231231231` are counted as links, but they are counted as broken(inaccessible) links as well

### Counted as broken but not really!
To identify that a link is broken or not, we rely on returned status code. Typically, a status code of 400 or above indicates a broken link. However, some websites do not respect this. For example, linked in replies to the request with the status code of 999. Try:
```bash
curl -v https://www.linkedin.com/company/home24/
````

This happens to many other websites, such as instagram, twitter and etc.

In cases like above, the link is counted as broken, although it is not really.

