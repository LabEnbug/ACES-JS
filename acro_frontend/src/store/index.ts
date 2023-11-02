import defaultSettings from '../settings.json';
import axios from 'axios';

export interface GlobalState {
  settings?: typeof defaultSettings;
  userInfo?: {
    username?: string;
    nickname?: string;

    permissions: Record<string, string[]>;
  };
  userLoading?: boolean;
  baxios?: any;
}

const initialState: GlobalState = {
  settings: defaultSettings,
  userInfo: {
    permissions: {},
  },
  baxios: null
};

export default function store(state = initialState, action) {
  switch (action.type) {
    case 'update-settings': {
      const { settings } = action.payload;
      return {
        ...state,
        settings,
      };
    }
    case 'update-userInfo': {
      const { userInfo = initialState.userInfo, userLoading } = action.payload;
      return {
        ...state,
        userLoading,
        userInfo,
      };
    }
    case 'update-baxios': {
      const { baxios } = action.payload;
      return {
        ...state,
        baxios,
      };
    }
    default:
      return state;
  }
}
