'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import { ImageList } from '../components/images.jsx'


export class ImageBox extends React.Component{
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      pollInterval: this.props.pollInterval,
      url: this.props.url,
    }
  };
  loadCommentsFromServer() {
    $.ajax({
      url: this.state.url,
      dataType: 'json',
      cache: true,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.state.url, status, err.toString());
      }.bind(this)
    })
  }
  componentDidMount() {
    this.loadCommentsFromServer();
    setInterval(this.loadCommentsFromServer.bind(this), this.state.pollInterval);
  }
  render() {
    return (
      <div className="container">
        <ImageList data={this.state.data} />
      </div>
    );
  }
}
