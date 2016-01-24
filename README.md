# Papyrus

Papyrus is a real-time collaborative Markdown editor and document repository with simple organization and project-based management. At [Furqan Software](https://furqansoftware.com) we always wanted a simple way to collaborate on Markdown documents and Papyrus is our stab at fulfilling that need.

![Obligatory GIF](http://i.imgur.com/WhQqbFA.gif)

As of writing this README.md, [GopherGala 2016](http://gophergala.com/)'s 48 hours is almost up. The core collaborative editing functionality, namely operational transformation and other relevant constructs, have been implemented. The webapp is lacking some functionalities (e.g. removing a member that has been added to a project, deleting a document, etc). At this moment, you can login using a Google or GitHub account, create organizations, create projects, create documents and start editing them collaboratively in real-time.

## Usage

Create a _GOPATH_ for this project:

```
$ mkdir -p Papyrus
$ cd Papyrus
$ export GOPATH=`pwd`
```

Clone this repository:

```
$ mkdir -p src/github.com/gophergala2016/papyrus
$ cd src/github.com/gophergala2016/papyrus
$ git clone https://github.com/gophergala2016/papyrus.git .
```

Install dependencies:

```
$ go get -v ./...
```

Build _papyrusd_ binary:

```
$ go install ./cmd/papyrusd
```

Set necessary environment variables (see _env-sample.txt_):

```
$ export ADDR=:5000
$ export SECRET=some-long-secret
$ export BASE=http://localhost:5000
$ export MONGO_URL=mongodb://localhost/papyrus
$ export GITHUB_CLIENT_ID=your-github-client-id-here
$ export GITHUB_CLIENT_SECRET=your-github-client-secret-here
$ export GOOGLE_CLIENT_ID=your-google-client-id-here
$ export GOOGLE_CLIENT_SECRET=your-google-client-secret-here
```

Run _papyrusd_:

```
$GOPATH/bin/papyrusd
```

Make sure MongoDB is running and all necessary environment variables are set.

## Acknowledgements

- [CodeMirror](https://codemirror.net/) - One fantastic text editor for the web
- [Code Commit](http://www.codecommit.com/blog/java/understanding-and-applying-operational-transformation) - For their amazing explanation of operational transformation
- [OT Explained](http://operational-transformation.github.io/visualization.html) - That visualization!
- [golab](https://github.com/mb0/lab) - For the insights we got from its code

## License

Loadcat is available under the [BSD (3-Clause) License](http://opensource.org/licenses/BSD-3-Clause).
