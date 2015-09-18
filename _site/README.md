# How to contribute

[*CLI web page*](http://cli.scalingo.com) relies on Jekyll and GitHub Pages to host, build and serve. You are welcome to contribute to this documentation by forking and sending pull requests.

## Adding a new changelog entry

To add a new article simply create a new markdown file in the `_posts` folder with the following file format: `yyyy-mm-dd-title.md` (*e.g. 1970-01-01-what-is-epoch.md*).
Please use the current date when creating the file.

The minimal *front matter* that you need to add is:

```
---
title: What is Epoch
modified_at: 1970-01-01 00:00:00
categories: unix
tags: time
---
```

Be aware that the `categories` will be used in the final URL. In the case of the previous *front matter*, the generated URL will be: `/unix/what-is-epoch.html`.

## Modifying an article

You are welcome to modify any article, but please remember to update `modified_at` before sending your pull request.

## Don'ts

Please do not use the `date` metadata as it will conflict with the date extracted from the file name. Instead, use `modified_at` to record when a modification is made to an article.

Please do not use first-level HTML/Markdown headers (*i.e. `<h1></h1>`*) as it will be pulled from the `title ` metadata.

## Running locally

To install dependancies locally:

```
bundle install
```

To build the static site and spin-up a file server:

```
bundle exec jekyll server
```

## Links

* [Scalingo Documentation Center](http://doc.scalingo.com)
* [Using Jekyll with Pages](https://help.github.com/articles/using-jekyll-with-pages/)
* [Jekyll](https://jekyllrb.com/)
