'use strict';

import React from 'react';
import { Header } from '../components/header.jsx'
import { Footer } from '../components/footer.jsx'
import { ImageBox } from '../boxes/images.jsx'

class FeaturedImages extends React.Component {
  render(){
    return (
      <div>
        <Header/>
        <ImageBox url="/api/v0/images/featured" pollInterval={10000000} />
        <Footer/>
      </div>
    )
  }
}

export { FeaturedImages }
