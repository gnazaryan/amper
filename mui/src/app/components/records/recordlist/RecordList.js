import React, { useState, useEffect, useImperativeHandle, useRef, forwardRef} from 'react';
import { DataGridPremium } from '../../x-data-grid-premium';

import DataStore from '../../../data/DataStore';
import HostManager from '../../../../HostManager';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import ReferenceRenderer from './cellrenderers/ReferenceRenderer';
import Typography from '@mui/material/Typography';
import { debounceLatest } from '../../../amper/Instruments';
import {booleanOnlyOperators} from './filters/BooleanFilter'
import {objectTypeOnlyOperators} from './filters/ObjectTypeFilter'
import { anyOperators } from './filters/AnyFilter';
import { referenceOperators } from './filters/ReferenceFilter';
import { numberOnlyOperators } from './filters/NumberFilter';
import { RecordToolbar } from './RecordToolbar';
import AmperConstatns from '../../../util/AmperConstants';
import { post } from '../../../data/Submit';
import BooleanRenderer from './cellrenderers/BooleanRenderer';
import BooleanEditRenderer from "./celleditrenderers/BooleanEditRenderer";
import TextEditRenderer from './celleditrenderers/TextEditRenderer';
import ReferenceEditRenderer from './celleditrenderers/ReferenceEditRenderer';
import {parseBoolean} from '../../../../app/util/BooleanUtil';
import NumberRenderer from './cellrenderers/NumberRenderer';
import NumberEditRenderer from './celleditrenderers/NumberEditRenderer';
import DateRenderer from './cellrenderers/DateRenderer';
import DateEditRenderer from './celleditrenderers/DateEditRenderer';
import clsx from 'clsx';
import DateTimeRenderer from './cellrenderers/DateTimeRenderer';
import DateTimeEditRenderer from './celleditrenderers/DateTimeEditRenderer';
import TextRenderer from './cellrenderers/TextRenderer';
import Convenience from '../../../help/Convenience';
import { AppContext } from '../../../../App';
import Popover from '@mui/material/Popover';

function RecordList(props, ref) {
    const app = React.useContext(AppContext);
    const {id, onStateChange, toast, onSelect, onLoad} = props;
    const existingOrder = useRef([]);
    const onColumnOrderChange = (order) => {
        if (existingOrder.current) {
            const idx = existingOrder.current.indexOf(order.column.field);
            existingOrder.current.splice(idx, idx !== -1 ? 1 : 0);
            existingOrder.current.splice(order.targetIndex - 1, 0, order.column.field);
            const newState = {
                ...state,
                state: {
                    ...state.state,
                    columnsOrder: existingOrder.current,
                }
            };
            setState(newState);
            if (onStateChange) {
                onStateChange({
                    ...state.state,
                    columnsOrder: existingOrder.current,
                }, true)
            }
        }
    };
    const onFilterChangeUpdate = (getFilterModel) => {
        if (onStateChange) {
            onStateChange({
                ...state.state,
                filterModel: getFilterModel(),
            }, false)
        }
    }
    const getFilterModelMap = (filterModel) => {
        const filterModelMap = {};
        if (filterModel.items) {
            for (let i = 0; i < filterModel.items.length; i++) {
                const filter = filterModel.items[i];
                if (filterModelMap[filter.columnField] == null) {
                    filterModelMap[filter.columnField] = {}
                }
                filterModelMap[filter.columnField][filter.id] = filter;
            }
        }
        return filterModelMap;
    };
    const refresh = () => {
        setState({
            ...state,
            loading: true,
            payloadsMap: {},
            totalCount: null,
            previousTotalCount: state.totalCount,
        });
        if (myRef.current.disableUpdate) {
            myRef.current.disableUpdate();
        }
    };
    
    const LIMIT = 100;
    const initialState = () => {
         return {
            firstTimeLoading: (props.object != null && props.object.apiName != null), 
            loading: (props.object != null && props.object.apiName != null),
            startId: 0,
            page: 0,
            data: [],
            metadata: null,
            totalCount: null,
            state: props.state || {},
            onFilterChangeDebounce: debounceLatest(onFilterChangeUpdate),
            filterModel: props.state.filterModel || {items: []},
            filterModelMap: getFilterModelMap(props.state.filterModel || {}),
            payloadsMap: {},
            cache: {},
            sort: props.state != null && props.state.sort != null ? props.state.sort : {},
            limit: props.state != null && props.state.limit != null ? props.state.limit : LIMIT,
            recordPaging: {
                page: 0,
                pageSize: 100,
            }
        };
    };
    const [state, setState] = useState(initialState);
    const [filterButtonEl, setFilterButtonEl] = React.useState(null);
    const myRef = useRef();
    if (app) {
        app.registerRefresh(id, () => {
            refresh();
        });    
    }
    useImperativeHandle(ref, () => ({
        reset() {
            setState(initialState);
        },
        addFilter(filter) {
            const filterModel = state.filterModel;
            const items = [filter];
            for (let i = 0; i < filterModel.items.length; i++) {
                //check if the filter is not a interaction filter, if the filter was added
                //as interaction then no need to add it
                if (filterModel.items[i].id < 100000000000000 && filterModel.items[i].id % 1 === 0) {
                    items.push(filterModel.items[i])
                }
            }
            filterModel.items = items;
            setState({
                ...state,
                filterModel,
                filterModelMap: getFilterModelMap(filterModel),
                loading: true,
            });
        }
    }));

    useEffect(() => {
        if (state.loading) {
          getDataStore().load((result)=> {
              setState({
                ...state,
                loading: false,
                data: result.data || [],
                metadata: state.metadata ? state.metadata : result.metadata,
                firstTimeLoading: false,
                totalCount: state.totalCount != null ? state.totalCount : result.totalCount,
              })
              onLoad(state.metadata ? state.metadata : result.metadata);
          });
        }
      }, [state.loading]);


    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}records/fetch`,
            requestMethod: "POST",
            parameters: {
                apiName: props.object.apiName,
                objectType: props.objectType ? props.objectType.apiName : null,
                start: (state.recordPaging.pageSize * state.recordPaging.page),
                limit: state.recordPaging.pageSize,
                search: getSearch(),
                metadata: state.metadata == null,
            }
        });
    };
    
    const getSearch = () => {
        const result = {
            totalCount: state.totalCount == null,
        };
        if (state.filterModel &&  state.filterModel.items && state.filterModel.items.length) {
            const filters = [];
            for (let i = 0; i < state.filterModel.items.length; i++) {
                const item = state.filterModel.items[i];
                filters.push({
                    apiName: item.columnField,
                    operator: item.operatorValue,
                    value: item.value,
                })
            }
            result.operator = state.filterModel.linkOperator;
            if (state.filterModel.quickFilterValues && state.filterModel.quickFilterValues.length > 0) {
                result.quickOperator = state.filterModel.quickFilterLogicOperator;
                result.quickValues = state.filterModel.quickFilterValues;
            }
            if (filters.length > 0) {
                result.term = filters
            }
        }
        if (state.sort != null && state.sort.field != null && state.sort.dir != null) {
            result.sortField = state.sort.field;
            result.sortDir = state.sort.dir;
        }
        return JSON.stringify(result);
    };

    const getProgressBar = () => {
        return <Box sx={{ display: 'flex', width: '100%', height: '100%', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
            <CircularProgress />
        </Box>;
    };

    const filterApply = () => {
        setState({
            ...state,
            loading: true,
            totalCount: null,
            previousTotalCount: state.totalCount,
        });
    };

    const getPayloadValue = (identifier, key) => {
        if (state.payloadsMap[identifier] != null && state.payloadsMap[identifier][key] != null) {
            return state.payloadsMap[identifier][key];
        }
        return null;
    };

    const getCacheReferenceLabel = (identifier, key) => {
        const record = state.cache[identifier];
        if (record != null) {
            const referenceRecord = record[key];
            if (referenceRecord != null && referenceRecord[AmperConstatns.SYSTEM_FIELDS.NAME] != null) {
                return referenceRecord[AmperConstatns.SYSTEM_FIELDS.NAME];
            }
        }
        return null;
    };

    const cacheValue = (identifier, key, value, originalValue) => {
        setState({
            ...state,
            cache: {
                ...state.cache,
                [identifier]: {
                    ...state.cache[identifier],
                    [key]: originalValue,
                }
            },
            payloadsMap: {
                ...state.payloadsMap,
                [identifier]: {
                    ...state.payloadsMap[identifier],
                    [key]: value,
                },
            },
        });
        if (identifier && key && value != null && myRef.current.enableUpdate) {
            myRef.current.enableUpdate()
        }
    };

    const getCacheValue = (identifier, key) => {
        if (state.cache[identifier] != null && state.cache[identifier][key] != null) {
            return state.cache[identifier][key];
        }
        return null;
    };

    const cellClassName = (params) => {
        if (params.value == null) {
          return '';
        }

        let dirty = false;
        if (params.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER] && 
            state.payloadsMap[params.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER]] != null &&
            state.payloadsMap[params.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER]][params.field] != null) {
          dirty = true;
        }
        return clsx('amper-records', {
          normal: !dirty,
          dirty: dirty,
        });
      };
    const widths = {
        id: 20,
        status_sys: 120
    }
    const getWidth = (key) => {
        if (state.state.columnsSize != null && state.state.columnsSize[key] != null) {
            return state.state.columnsSize[key];
        }
        return widths[key] == null ? 200 : widths[key];
    };
    const NON_EDITABLE_FIELDS = ['id', 'identifier_sys', 'objectType_sys', ]
    const getColumns = () => {
        const result = [];
        existingOrder.current = []
        if (!state.metadata) {
            return result;
        }
        const columnsMap = {};
        for (let i = 0; i < state.metadata['Fields'].length; i++) {
            let field = state.metadata['Fields'][i];
            columnsMap[field.apiName] = field;
        }
        const columnsOrder = state.state.columnsOrder || ['name_sys', 'objectType_sys', 'status_sys']
        for (let l = 0; l < columnsOrder.length; l++) {
            const field = columnsMap[columnsOrder[l]]
            if (field != null) {
                const column = {
                    field: field.apiName,
                    headerName: field.apiName,
                    headerNameLabel: field.label,
                    width: getWidth(field.apiName),
                    editable: !NON_EDITABLE_FIELDS.includes(field.apiName),
                    required: parseBoolean(field.required),
                    getPayloadValue,
                    cacheValue,
                    getCacheValue,
                    cellClassName: cellClassName,
                    renderHeader: (params) => {
                        return params.colDef.headerNameLabel;
                    }
                };
                if (field.type === 'REFERENCE') {
                    column.renderCell = ReferenceRenderer;
                    column.renderEditCell = ReferenceEditRenderer;
                    column.referenceObjectId = field.objectReference;
                    column.filterOperators = referenceOperators(filterApply, field.objectReference);
                } else if (field.type === 'BOOLEAN') {
                    column.filterOperators = booleanOnlyOperators(filterApply);
                    column.renderCell = BooleanRenderer;
                    column.renderEditCell = BooleanEditRenderer;
                } else if (field.apiName === 'objectType_sys') {
                    column.filterOperators = objectTypeOnlyOperators(filterApply, state.metadata['ObjectTypes'] || []);
                } else if (field.type === 'INTEGER' || field.type === 'NUMBER') {
                    column.filterOperators = numberOnlyOperators(filterApply);
                    column.maxLength = field.textLength;
                    column.renderCell = NumberRenderer;
                    column.renderEditCell = NumberEditRenderer;
                } else if (field.type === 'TEXT') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = TextRenderer;
                    column.renderEditCell = TextEditRenderer;
                } else if (field.type === 'DATE') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = DateRenderer;
                    column.renderEditCell = DateEditRenderer;
                } else if (field.type === 'DATETIME') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = DateTimeRenderer;
                    column.renderEditCell = DateTimeEditRenderer;
                } else {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = TextRenderer;
                    column.renderEditCell = TextEditRenderer;
                }
                result.push(column)
                existingOrder.current.push(field.apiName);
                delete columnsMap[columnsOrder[l]]
            }
        }
        for (const [apiName, field] of Object.entries(columnsMap)) {
            if (field != null) {
                const column = {
                    field: field.apiName,
                    headerName: field.apiName,
                    headerNameLabel: field.label,
                    width: getWidth(field.apiName),
                    editable: !NON_EDITABLE_FIELDS.includes(field.apiName),
                    required: parseBoolean(field.required),
                    getPayloadValue,
                    cacheValue,
                    getCacheValue,
                    cellClassName: cellClassName,
                    renderHeader: (params) => {
                        return params.colDef.headerNameLabel;
                    }
                };
                if (field.type === 'REFERENCE') {
                    column.renderCell = ReferenceRenderer;
                    column.renderEditCell = ReferenceEditRenderer;
                    column.referenceObjectId = field.objectReference;
                    column.filterOperators = referenceOperators(filterApply, field.objectReference);
                } else if (field.type === 'BOOLEAN') {
                    column.renderCell = BooleanRenderer;
                    column.renderEditCell = BooleanEditRenderer;
                    column.filterOperators = booleanOnlyOperators(filterApply);
                } else if (field.apiName === 'objectType_sys') {
                    column.filterOperators = objectTypeOnlyOperators(filterApply, state.metadata['ObjectTypes'] || []);
                } else if (field.type === 'INTEGER' || field.type === 'NUMBER') {
                    column.filterOperators = numberOnlyOperators(filterApply);
                    column.maxLength = field.textLength;
                    column.renderCell = NumberRenderer;
                    column.renderEditCell = NumberEditRenderer;
                } else if (field.type === 'TEXT') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = TextRenderer;
                    column.renderEditCell = TextEditRenderer;
                } else if (field.type === 'DATE') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = DateRenderer;
                    column.renderEditCell = DateEditRenderer;
                } else if (field.type === 'DATETIME') {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = DateTimeRenderer;
                    column.renderEditCell = DateTimeEditRenderer;
                } else {
                    column.filterOperators = anyOperators(filterApply);
                    column.renderCell = TextRenderer;
                    column.renderEditCell = TextEditRenderer;
                }
                existingOrder.current.push(field.apiName);
                result.push(column)
            }
        }
        return result;
    };

    const setPage = (newPage) => {
        let startId = state.startId;
        if (state.data && state.data.length > 0) {
            if (newPage > state.page) {
                startId = state.data[state.data.length - 1]['id']
            } else {
                startId = state.previousStartIds[newPage]
            }
        }
        startId = parseInt(startId)

        setState({
            ...state,
            previousStartIds: {
                ...state.previousStartIds,
                [state.page]: state.startId,
            },
            page: newPage,
            startId: startId,
            loading: startId !== state.startId,
        });
    };

    const onColumnOrderChangeDebounceCallee = (order) => {
        onColumnOrderChange(order)
    };

    const onColumnVisibilityChangeDebounceCallee = (columnVisibilityModel) => {
        const result = [];
        for (const [key, value] of Object.entries(columnVisibilityModel)) {
            if (value === false) {
                result.push(key)
            }
        }
        if (onStateChange) {
            onStateChange({
                ...state.state,
                hiddenColumns: result,
            }, false)
        }
    }

    const getHiddenColumns = () => {
        const result = {};
        if (state.state && state.state.hiddenColumns) {
            for (let i = 0; i < state.state.hiddenColumns.length; i++) {
                result[state.state.hiddenColumns[i]] = false;
            }
        }
        return result;
    };

    const onFilterChange = (newFilterModel, action) => {
        const filterModelMap = {};
        let update = false;
        for (let i = 0; i < newFilterModel.items.length; i++) {
            const filter = newFilterModel.items[i];
            if (filterModelMap[filter.columnField] == null) {
                filterModelMap[filter.columnField] = {}
            }
            filterModelMap[filter.columnField][filter.id] = filter;

            //Check if a filter operator was changed
            if (state.filterModelMap && state.filterModelMap[filter.columnField] != null
                 && state.filterModelMap[filter.columnField][filter.id] != null
                 && state.filterModelMap[filter.columnField][filter.id].operatorValue !== filter.operatorValue) {
                update = true;
            }
        }

        //Check if a filter was removed
        if (!update && state.filterModel && state.filterModel.items.length > newFilterModel.items.length) {
            update = true;
        }

        //Check if operator is updated
        if (!update && state.filterModel && state.filterModel.linkOperator != newFilterModel.linkOperator) {
            update = true;
        }
        setState(previousState => {
            previousState.filterModelMap = filterModelMap;
            previousState.filterModel = newFilterModel;
            return {
                ...previousState,
                loading: update,
                totalCount: update ? null : previousState.totalCount,
                previousTotalCount: state.totalCount,
            };
        });
        state.onFilterChangeDebounce(newFilterModel);
    };

    const onRemove = () => {
        if (state.selectedRowData && state.selectedRowData.length > 0) {
            const identifiers = state.selectedRowData.filter(record => record[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER] != null).map(record => record[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER])
            if (identifiers && identifiers.length > 0) {
                post(`${HostManager.amperHost()}records/removeRecords`, {
                    identifiers: identifiers
                  }, (result) => {
                    app.toast('info', `The object "${state.metadata.Object.apiName}" records were successfully removed.`);
                    refresh();
                  }, (result) => {
                    app.toast('warning', result.error != null ? result .error : `The object "${state.metadata.Object.apiName}" records removal was not successfull`);
                  });
            }
        }
    };

    const onUpdate = () => {
        const payloads = [];
        for (const [key, value] of Object.entries(state.payloadsMap)) {
            payloads.push({
                ...value,
                [AmperConstatns.SYSTEM_FIELDS.IDENTIFIER]: key,
            })
        }
        if (payloads.length > 0) {
            post(`${HostManager.amperHost()}records/updateRecords`, {
                payloads: JSON.stringify(payloads)
              }, (result) => {
                app.toast('info', 'The records were successfully updated.');
                refresh();
              }, (result) => {
                app.toast('warning', 'The records were not all successfully updated.');
              });
        }
    };

    const onPageSizeChange = (pageSize) => {
        setState({
            ...state,
            loading: true,
            limit: pageSize,
        });
        if (onStateChange) {
            onStateChange({
                ...state.state,
                limit: pageSize,
            }, false);
        }
    };

    const onColumnWidthChange = (params) => {

        setState({
            ...state,
            state: {
                ...state.state,
                columnsSize: {
                    ...state.state.columnsSize,
                    [params.colDef.field]: params.colDef.width,
                },
            }
        });
        if (onStateChange) {
            onStateChange({
                ...state.state,
                columnsSize: {
                    ...state.state.columnsSize,
                    [params.colDef.field]: params.colDef.width,
                },
            }, false);
        }
    };

    const isCellEditable = (params) => {
        let result = false;
        const objectTypeName = params.row[AmperConstatns.SYSTEM_FIELDS.OBJECT_TYPE]
        if (Convenience.hasValue(objectTypeName)) {
            for (let i = 0; i < state.metadata.ObjectTypes.length; i++) {
                const objectType = state.metadata.ObjectTypes[i];
                if (objectTypeName === objectType.apiName) {
                    for (let l = 0; l < objectType.objectTypeFields.length; l++) {
                        const objectTypeField = objectType.objectTypeFields[l];
                        if (params.field === objectTypeField.apiName) {
                            return true;
                        }
                    }
                }
            }
        }


        return result;
    };

    const setRecordsPaginationModel = (pagingModel) => {
        setState({
            ...state,
            recordPaging: pagingModel,
            loading: true,
        });
    };

    const onSortModelChange = (model) => {
        let field = null;
        let dir = null;
        if (model.length > 0) {
            field = model[0].field;
            dir = model[0].sort;
        }

        setState({
            ...state,
            sort: {
                field: field,
                dir: dir
            },
            loading: true,
        });
        if (onStateChange) {
            onStateChange({
                ...state.state,
                sort: {
                    field: field,
                    dir: dir
                },
            }, false);
        }
    };

    const getSortModel = () => {
        const result =[];
        if (state.sort && state.sort.field) {
            result.push({
                field: state.sort.field,
                sort: state.sort.dir
            });
        }
        return result;
    };

    const  isOverflown = (element) => {
        return element.scrollHeight > element.clientHeight || element.scrollWidth > element.clientWidth;
    };

    const handlePopoverOpen = (event, arg, arg1) => {
        if (isOverflown(event.currentTarget)) {
            const field = event.currentTarget.dataset.field;
            if (AmperConstatns.SYSTEM_FIELDS.STATUS === field) {
                return;
            }
            const id = event.currentTarget.parentElement.dataset.id;
            const row = state.data.find((r) => r.id === id);
            const referenceValue = row[field + '_name_sys'];
            //First attempt to show the cached value
            //If not found implies field is not modified
            //then try to show the original value
            let cachedValue = getCacheReferenceLabel(row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], field);
            if (cachedValue == null) {
                cachedValue = getPayloadValue(row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], field);
            }
            if ((cachedValue == null || cachedValue == '') && (referenceValue == null || referenceValue =='') && (row[field] == null || row[field] == '') ) {
                return;
            }
            setValue(cachedValue || referenceValue || row[field]);
            setPopoverAnchorEl(event.currentTarget);    
        }
        event.stopPropagation();
    };

    const handlePopoverClose = () => {
        setPopoverAnchorEl(null);
    };
    const [popoverAnchorEl, setPopoverAnchorEl] = React.useState(null);
    const [value, setValue] = React.useState('');

    const getRecordListGrid = () => {
        const hiddenColumns = getHiddenColumns();
        const columns = getColumns();
        return <Box key={id}
                    sx={{
                    height: '100%',
                    width: '100%',
                    '& .amper-records.normal': {
                        fontSize: '1rem',
                    },
                    '& .amper-records.dirty': {
                        fontSize: '1rem',
                        backgroundColor: 'inactive.main',
                        color: 'primary.contrastText',
                    },
                    }}
                >
            <DataGridPremium
                initialState={{
                    columns: {
                    columnVisibilityModel: hiddenColumns,
                    },
                }}
                isCellEditable={isCellEditable}
                onRowSelectionModelChange={(ids) => {
                    const selectedIDs = new Set(ids);
                    const rowData = state.data.filter((row) =>
                        selectedIDs.has(row.id)
                    )

                    setState({
                        ...state,
                        selectedRowData: rowData,
                    });
                    setTimeout(() => {
                        onSelect(id, state.metadata, rowData);
                    }, 1);
                }}
                filterMode='server'
                onFilterModelChange={onFilterChange}
                filterModel={state.filterModel}
                rowCount={state.totalCount != null ? state.totalCount : state.previousTotalCount}
                paginationMode="server"
                rows={state.data}
                slots={{
                    toolbar: RecordToolbar,
                }}
                slotProps={{
                    panel: {
                        anchorEl: filterButtonEl,
                    },
                    toolbar: {
                        setFilterButtonEl,
                        metadata: state.metadata,
                        toast,
                        refresh,
                        onRemove,
                        onUpdate,
                        parentRef: myRef,
                    },
                    cell: {
                        onMouseEnter: handlePopoverOpen,
                        onMouseLeave: handlePopoverClose,
                      },
                }}
                columns={columns}
                loading={state.loading}
                onPageSizeChange={onPageSizeChange}
                onSortModelChange={onSortModelChange}
                sortModel={getSortModel()}
                onColumnWidthChange={onColumnWidthChange}
                pageSizeOptions={[100, 500, 1000, 5000, 50000]}
                onPaginationModelChange={setRecordsPaginationModel}
                paginationModel={state.recordPaging}
                checkboxSelection
                onColumnOrderChange={onColumnOrderChangeDebounceCallee}
                onColumnVisibilityModelChange={onColumnVisibilityChangeDebounceCallee}
                pagination
                disableSelectionOnClick
                experimentalFeatures={{ newEditingApi: true }}
                onPageChange={(newPage) => setPage(newPage)}
            />
            <Popover
                sx={{ pointerEvents: 'none'}}
                open={Boolean(popoverAnchorEl)}
                anchorEl={popoverAnchorEl}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                onClose={handlePopoverClose}
                disableAutoFocus
                disableRestoreFocus>
                <Typography sx={{ p: 1, maxWidth: 500 }}>{value}</Typography>
            </Popover>
        </Box>;
    };

    const getNotConfigured = () => {
        return <Box sx={{ display: 'flex', width: '100%', height: '100%', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
            <Typography variant="subtitle1" gutterBottom>
                Record list is not configured.
            </Typography>
        </Box>;
    };
    return state.firstTimeLoading ? getProgressBar() : (props.object && props.object.apiName) ? getRecordListGrid() : getNotConfigured();
};
export default forwardRef(RecordList);