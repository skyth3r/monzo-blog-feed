# Monzo Blog Feed

A Go script that generates rss+xml, atom & json feeds for the [Monzo blog](https://monzo.com/blog) that can be used by rss reader apps.

I put together this script to keep up to date with new posts from Monzo via my rss reader apps (Currently I use [Tapestry](https://usetapestry.com/), and [NetNewsWire](https://netnewswire.com/)).

[Monzo switched blogging platforms over 5 years ago](https://twitter.com/monzo/status/1294232587280568322), and the blogging platform doesn't seem to support rss feeds...so I generated the feeds myself with a little bit of concurrent web scrapping.

## Feeds

You can use these feeds below in any rss reader app you use to keep up to date with new posts on the Monzo blog.

There are three feeds available:

- Monzo blog - all posts from the main blog
- Monzo blog - Technology - All posts tagged with the 'Technology' tag
- Monzo US blog - posts from the US blog

These feeds are updated automatically using a GitHub Actions workflow.

### RSS feeds

- [Monzo Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog.rss)
- [Monzo Blog - Technology](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog_technology.rss)
- [Monzo US Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/us_blog.rss)

### Atom feeds
- [Monzo Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog.atom)
- [Monzo Blog - Technology](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog_technology.atom)
- [Monzo US Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/us_blog.atom)

### JSON feeds

- [Monzo Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog.json)
- [Monzo Blog - Technology](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/blog_technology.json)
- [Monzo US Blog](https://raw.githubusercontent.com/skyth3r/monzo-blog-feed/refs/heads/main/feeds/us_blog.json)