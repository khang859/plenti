package generated

// Do not edit, this file is automatically generated.

// Ejected: scaffolding used in 'build' command
var Ejected = map[string][]byte{
	"/build.js": []byte(`import svelte from 'svelte/compiler.js';
import 'svelte/register.js';
import relative from 'require-relative';
import path from 'path';
import fs from 'fs';

// Get the arguments from Go command execution.
const args = process.argv.slice(2)

// -----------------
// Helper Functions:
// -----------------

// Create any missing sub folders.
const ensureDirExists = filePath => {
	let dirname = path.dirname(filePath);
	if (fs.existsSync(dirname)) {
		return true;
	}
	ensureDirExists(dirname);
	fs.mkdirSync(dirname);
}

// Concatenates HTML strings together.
const injectString = (order, content, element, html) => {
	if (order == 'prepend') {
		return html.replace(element, content + element);
	} else if (order == 'append') {
		return html.replace(element, element + content);
	}
};

// -----------------------
// Start client SPA build:
// -----------------------

let clientBuildStr = JSON.parse(args[0]);

clientBuildStr.forEach(arg => {

	let layoutPath = path.join(path.resolve(), arg.layoutPath)
	let component = fs.readFileSync(layoutPath, 'utf8');

	// Create component JS that can run in the browser.
	let { js, css } = svelte.compile(component, {
		css: false
	});
	  
	// Write JS to build directory.
	ensureDirExists(arg.destPath);
	fs.promises.writeFile(arg.destPath, js.code);

	// Write CSS to build directory.
	ensureDirExists(arg.stylePath);
	if (css.code && css.code != 'null') {
		fs.appendFileSync(arg.stylePath, css.code);
	}
});

// ------------------------
// Start static HTML build:
// ------------------------

let staticBuildStr = JSON.parse(args[1]);
let allNodes = JSON.parse(args[2]);

// Create the component that wraps all nodes.
let htmlWrapper = path.join(path.resolve(), 'layout/global/html.svelte')
const component = relative(htmlWrapper, process.cwd()).default;

staticBuildStr.forEach(arg => {

	let componentPath = path.join(path.resolve(), arg.componentPath);
	let destPath = path.join(path.resolve(), arg.destPath);

	// Set route used in svelte:component as "this" value.
	const route = relative(componentPath, process.cwd()).default;

	// Set props so component can access field values, etc.
	let props = {
		route: route,
		node: arg.node,
		allNodes: allNodes
	};

	// Create the static HTML and CSS.
	let { html, css } = component.render(props);

	// Inject Style.
	let style = "<style>" + css.code + "</style>";
	html = injectString('prepend', style, '</head>', html);
	// Inject SPA entry point.
	let entryPoint = '<script type="module" src="https://unpkg.com/dimport?module" data-main="/spa/ejected/main.js"></script><script nomodule src="https://unpkg.com/dimport/nomodule" data-main="/spa/ejected/main.js"></script>';
	html = injectString('prepend', entryPoint, '</head>', html);
	// Inject ID used to hydrate SPA.
	let hydrator = ' id="hydrate-plenti"';
	html = injectString('append', hydrator, '<html', html);

	// Write .html file to filesystem.
  	ensureDirExists(destPath);
	fs.promises.writeFile(destPath, html);
	  
});`),
	"/main.js": []byte(`import Router from './router.svelte';

/*
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/plenti-service-worker.js')
  .then((reg) => {
    console.log('Service Worker registration succeeded.');
  }).catch((error) => {
    console.log('Service Worker registration failed with ' + error);
  });
} else {
  console.log('Service Workers not supported by browser')
}
*/

const replaceContainer = function ( Component, options ) {
  const frag = document.createDocumentFragment();
  const component = new Component( Object.assign( {}, options, { target: frag } ));
  if (options.target) {
    options.target.replaceWith( frag );
  }
  return component;
}

const app = replaceContainer( Router, {
  target: document.querySelector( '#hydrate-plenti' ),
  props: {}
});

export default app;
`),
	"/router.svelte": []byte(`<Html {route} {node} {allNodes} />

<script>
  import Navaid from 'navaid';
  import nodes from './nodes.js';
  import Html from '../global/html.svelte';

  let route, node, allNodes;

  const getNode = (uri, trailingSlash = "") => {
    return nodes.find(node => node.path + trailingSlash == uri);
  }

  let uri = location.pathname;
  node = getNode(uri);
  if (node === undefined) {
    node = getNode(uri, "/");
  }
  allNodes = nodes;

  function draw(m) {
    node = getNode(uri);
    if (node === undefined) {
      // Check if there is a 404 data source.
      node = getNode("/404");
      if (node === undefined) {
        // If no 404.json data source exists, pass placeholder values.
        node = {
          "path": "/404",
          "type": "404",
          "filename": "404.json",
          "fields": {}
        }
      }
    }
    route = m.default;
    window.scrollTo(0, 0);
  }

  function track(obj) {
    uri = obj.state || obj.uri;
  }

  addEventListener('replacestate', track);
  addEventListener('pushstate', track);
  addEventListener('popstate', track);

  const handle404 = () => {
    import('../content/404.js')
      .then(draw)
      .catch(err => {
        console.log("Add a '/layout/content/404.svelte' file to handle Page Not Found errors.");
        console.log("If you want to pass data to your 404 component, you can also add a '/content/404.json' file.");
        console.log(err);
      });
  }

  const router = Navaid('/', handle404);

  allNodes.forEach(node => {
    router.on(node.path, () => {
      // Check if the url visited ends in a trailing slash (besides the homepage).
      if (uri.length > 1 && uri.slice(-1) == "/") {
        // Redirect to the same path without the trailing slash.
        router.route(node.path, false);
      } else {
        import('../content/' + node.type + '.js').then(draw).catch(handle404);
      }
    });

  });

  router.listen();

</script>
`),
}