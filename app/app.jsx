'use strict';

import ReactDOM from 'react-dom'
import React from 'react';
import { Router, Route, Link, browserHistory } from 'react-router'

import { FeaturedImages, UploadImage, ViewImage } from './views/images.jsx'
import { Signin, Signup, CreateUser, ViewUser, FeaturedUsers } from './views/user.jsx'
import { FeaturedCollections, ViewCollection } from './views/collections.jsx'

ReactDOM.render(
  <Router history={browserHistory}>
    <Route path="/" component={FeaturedImages}>
      <Route path="/upload" component={UploadImage}>
      </Route>

      <Route path="/signup" component={CreateUser}>
      </Route>

      <Route path="/signin" component={Signin}>
      </Route>

      <Route path="/collections" component={FeaturedCollections}>
        <Route path="/collections/:collectID" component={ViewCollection}>
        </Route>
      </Route>

      <Route path="/users" component={FeaturedUsers}>
        <Route path="/users/:userID" component={ViewUser}>
        </Route>
      </Route>

      <Route path="/images" component={FeaturedImages}>
        <Route path="/images/:imageID" component={ViewImage}>
        </Route>
      </Route>

    </Route>
  </Router>,
  document.getElementById('content')
);
