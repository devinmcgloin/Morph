'use strict';

import { FeaturedImages } from './views/featured.jsx'
import ReactDOM from 'react-dom'
import React from 'react';
import { Router, Route, Link, browserHistory } from 'react-router'

ReactDOM.render(
  <Router history={browserHistory}>
    <Route path="/" component={FeaturedImages}>
    </Route>
  </Router>,
  document.getElementById('content')
);
