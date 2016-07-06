'use strict';

import React from 'react';

export class Image extends React.Component{
  render() {
    return (
      <div className="image">
        <h2 className="image-src">
          <img src={this.props.url} />
        </h2>
      </div>
    );
  }
}


export class ImageList extends React.Component{
  render() {
    var commentNodes = this.props.data.map(function(img) {
      return (
        <Image url={img.source_link.medium} key={img.shortcode}>
          {img.shortcode}
        </Image>
      );
    });
    return (
      <div className="image-list">
        {commentNodes}
      </div>
    );
  }
}
