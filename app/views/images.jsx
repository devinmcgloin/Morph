'use strict';

import React from 'react';
import { Header } from '../components/header.jsx'
import { Footer } from '../components/footer.jsx'
import { ImageListBox, ImageBox } from '../boxes/images.jsx'

class FeaturedImages extends React.Component {
  render(){
    return (
        // <Header/>
          <ImageListBox url="/api/v0/images/featured" pollInterval={10000000} />
        // <Footer/>
    )
  }
}

class ViewImage extends React.Component {
  render(){
    return (
      // <Header/>
      <ImageBox url="/api/v0/images/{this.props.params.imageID}" pollInterval={10000000} />
      // <Footer/>
    )
  }
}

export { FeaturedImages }
