import React from 'react'

import ISkan  from './ISkan'

const ISkanViewer = ({ data }) => {
  console.log("iskan data", data);
  return <ISkan data={data} />
}

export default ISkanViewer;