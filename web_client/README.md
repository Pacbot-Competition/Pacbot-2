To compile this client, you will need Node installed. 

Steps to get a development server (for displaying code updates in real-time) for internal testing:
* `npm install` (only needs to be done once, after cloning)
* `npm run dev` (or `npm run host` for non-localhost access)

Steps to compile to a static HTML file for use during the competition:
* `npm install` (only needs to be done once, after cloning)
* `npm run prod`
* Navigate to `./dist/index.html` and open it in your default browser

Note `npm run host` and `npm run prod` are custom flags I introduced for this project - do not expect these commands to work in another Svelte project without updating the `package.json` file accordingly and copying the `./postProcess.js` script.