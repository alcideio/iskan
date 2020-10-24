import { forwardRef } from 'react';
import React from 'react';

import MaterialTable from 'material-table';
import { ThemeProvider as MuiThemeProvider } from '@material-ui/core/styles';
import { Typography } from '@material-ui/core'
import AddBox from '@material-ui/icons/AddBox';
import ArrowDownward from '@material-ui/icons/ArrowDownward';
import Check from '@material-ui/icons/Check';
import CheckCircle from '@material-ui/icons/CheckCircle';
import green from '@material-ui/core/colors/green';
import red from '@material-ui/core/colors/red';
import blue from '@material-ui/core/colors/blue';
import HighlightOff from '@material-ui/icons/HighlightOff';
import ChevronLeft from '@material-ui/icons/ChevronLeft';
import ChevronRight from '@material-ui/icons/ChevronRight';
import Clear from '@material-ui/icons/Clear';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import Edit from '@material-ui/icons/Edit';
import FilterList from '@material-ui/icons/FilterList';
import FirstPage from '@material-ui/icons/FirstPage';
import LastPage from '@material-ui/icons/LastPage';
import Remove from '@material-ui/icons/Remove';
import SaveAlt from '@material-ui/icons/SaveAlt';
import Search from '@material-ui/icons/Search';
import ViewColumn from '@material-ui/icons/ViewColumn';
import { createMuiTheme } from '@material-ui/core';

const tableIcons = {
    Add: forwardRef((props, ref) => <AddBox {...props} ref={ref} />),
    Check: forwardRef((props, ref) => <Check {...props} ref={ref} />),
    Clear: forwardRef((props, ref) => <Clear {...props} ref={ref} />),
    Delete: forwardRef((props, ref) => <DeleteOutline {...props} ref={ref} />),
    DetailPanel: forwardRef((props, ref) => <ChevronRight {...props} ref={ref} />),
    Edit: forwardRef((props, ref) => <Edit {...props} ref={ref} />),
    Export: forwardRef((props, ref) => <SaveAlt {...props} ref={ref} />),
    Filter: forwardRef((props, ref) => <FilterList {...props} ref={ref} />),
    FirstPage: forwardRef((props, ref) => <FirstPage {...props} ref={ref} />),
    LastPage: forwardRef((props, ref) => <LastPage {...props} ref={ref} />),
    NextPage: forwardRef((props, ref) => <ChevronRight {...props} ref={ref} />),
    PreviousPage: forwardRef((props, ref) => <ChevronLeft {...props} ref={ref} />),
    ResetSearch: forwardRef((props, ref) => <Clear {...props} ref={ref} />),
    Search: forwardRef((props, ref) => <Search {...props} ref={ref} />),
    SortArrow: forwardRef((props, ref) => <ArrowDownward {...props} ref={ref} />),
    ThirdStateCheck: forwardRef((props, ref) => <Remove {...props} ref={ref} />),
    ViewColumn: forwardRef((props, ref) => <ViewColumn {...props} ref={ref} />)
};

/*
          - Id: CVE-2017-16997
            Info:
              AffectedPackage: glibc
              AffectedVersion: 2.19-18+deb8u10
              CvssScore: "9.3"
              EffectiveSeverity: HIGH
              FixAvailable: "false"
              FixedPackage: glibc
              FixedVersion: ""
              Image: gcr.io/dcvisor-162009/alcide/dcvisor/cp-kafka@sha256:dabfd8697252225110d3bc8038d42219db498f1971eb514752030c6b7ded5975
            Severity: CRITICAL
*/


export default function CveTable({ findings }) {
    const theme = createMuiTheme({
        palette: {
          primary: {
            main: blue[500],
          },
          secondary: {
            main: green[500],
          },
        },
      });

    return (
    <MuiThemeProvider theme={theme}>
        <MaterialTable
        icons={tableIcons}
        title={<Typography variant="h4" component="h2" color="textPrimary">Image Vulnerabilities</Typography>}
        columns={[
            { 
                title: 'Fix', 
                field: 'Info.FixAvailable', 
                width: "10px",
                render: rowData => {
                    if (rowData.Info.FixAvailable === "true") {
                        return <CheckCircle style={{ color: green[500] }}/>
                    } else {
                        return <HighlightOff style={{ color: red[500] }}/>
                    }
                }
            },      
            { title: 'Severity', field: 'Severity', defaultSort: 'asc', width: "10px" },
            { title: 'Effective', field: 'Info.EffectiveSeverity', width: "10px" },            
            { title: 'CVSS', field: 'Info.CvssScore', width: "10px" },      
            { title: 'Package', field: 'Info.AffectedPackage', width: "10px" },
            { title: 'CVE', field: 'Id', width: "150px" },
            { title: 'Description', field: 'Info.Description' },
            { title: 'Image', field: 'Info.Image', defaultGroupOrder: 0 },            
        ]}
        data={findings}
        options={{
            padding: "dense",
            maxBodyHeight: 400,
            sorting: true,
            grouping: true,
            search: true,
            exportButton: true,
            //filtering: true
            headerStyle: {
                //backgroundColor: '#488ed8',
                color: blue[500],
                position: 'sticky'
            },
            rowStyle: (rowData) => {
                return {
                    borderBottom: "1px solid #cfd8dc",
                    borderCollapse: "collapse",
                }
            },
            paging: false
        }
        }
    />
    </MuiThemeProvider>)

}