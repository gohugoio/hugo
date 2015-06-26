---
aliases:
- /doc/localfiles/
date: 2015-06-12
menu:
  main:
    parent: extras
next: /community/mailing-list
notoc: true
prev: /extras/urls
title: Traversing Local Files
weight: 110
---

## Traversing Local Files

Hugo includes a way to traverse local files. 
This is done using the 'ReadDir' function. 

## Using ReadDir

ReadDir takes a single string input that is relative to the root directory of the site. It returns an array of [os.FileInfo](https://golang.org/pkg/os/#FileInfo)

Let's create a shortcode to build a file index with links using ReadDir. 

'fileindex.html'

    <table style="width=100%">
    <th>Size in bytes</th>
    <th>Name</th>
    {{$dir := .Get "dir"}}
    {{ $url := .Get "baseurl" }}
    
    {{ $files := ReadDir $dir }}
        {{ range $files }}
    			<tr>
                    <td>{{.Size}}</td>
                    <td>
                        <a href="{{$url}}{{.Name | urlize }}"> {{.Name}}</a>
                        </td>
                </tr> 
    	 {{ end }}
    </table>
    
Now lets use it to list the css files used on this site

    {{</* fileindex dir="static/css" baseurl="/css/" */>}}

Is rendered as:

{{< fileindex dir="static/css/" baseurl="/css/">}}
