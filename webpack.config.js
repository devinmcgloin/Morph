var webpack = require('webpack');

module.exports = {
     entry: './app/app.jsx',
     output: {
         path: './dist',
         filename: 'bundle.js',
     },
     module: {
         loaders: [{
             test: /\.jsx?$/,
             exclude: /node_modules/,
             loader: 'babel',
             query:
             {
               presets: ['es2015', 'react']
             }
         }]
     }
 }
