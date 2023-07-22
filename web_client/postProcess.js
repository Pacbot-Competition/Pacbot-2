import replace from 'replace';

/*  
  When vite compiles to a static dist/index.html file, 
  it generates a JS module, which crashes due to CORS 
  restrictions. This is a workaround I found - by 
  defering the JS script, it runs normally.
*/

replace({
  regex: 'type="module" crossorigin',
  replacement: 'defer',
  paths: ['./dist/index.html'],
  recursive: false,
  silent: false,
})