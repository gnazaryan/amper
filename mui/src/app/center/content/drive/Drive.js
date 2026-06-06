import React from 'react';
import { breadcrumbs } from '../../Breadcrambs'
import TileView from '../TileView';


export default function Drive({expanded}) {
  return (
        <TileView expanded={expanded} items={breadcrumbs.drive}></TileView>
    );
}
