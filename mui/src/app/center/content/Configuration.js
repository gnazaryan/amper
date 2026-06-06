import React from 'react';
import { breadcrumbs } from '../Breadcrambs'
import TileView from './TileView';


export default function Configuration({expanded}) {
  return (
        <TileView expanded={expanded} items={breadcrumbs.configuration}></TileView>
    );
}
