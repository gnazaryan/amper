import React from 'react';
import { breadcrumbs } from '../Breadcrambs'
import TileView from './TileView';


export default function Administration({expanded}) {
  return (
        <TileView expanded={expanded} items={breadcrumbs.administration}></TileView>
    );
}
