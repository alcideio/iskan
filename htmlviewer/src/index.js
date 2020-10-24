import React from 'react';
import ReactDOM from 'react-dom';
import 'regenerator-runtime';
import ISkanViewer from './ISkanViewer';

import './SkanViewer.module.scss';

var data = window['iskanReportData'];

ReactDOM.render(
  <React.StrictMode>
    <ISkanViewer data={data}/>
  </React.StrictMode>,
  document.getElementById('root')
);
