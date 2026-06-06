export const APP_ACTIONS = {
    EXPAND: 'EXPAND',
    COLAPSE: 'COLAPSE',
}

export default function AppReducer(state, action) {
  switch (action.type) {
    case APP_ACTIONS.EXPAND:
      return {
          ...state,
          expanded: true,
      };
    case APP_ACTIONS.COLAPSE:
      return {
          ...state,
          expanded: false,
      };
  }
};
