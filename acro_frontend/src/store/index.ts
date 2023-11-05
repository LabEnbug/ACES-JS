import defaultSettings from '../settings.json';
import axios from 'axios';

export interface GlobalState {
  settings?: typeof defaultSettings;
  isLogin?: boolean;
  userInfo?: {
    username?: string;
    nickname?: string;
    avatar_url?: string;
    balance?: number;

    permissions: Record<string, string[]>;
  };
  userLoading?: boolean;
  baxios?: any;
}

const initialState: GlobalState = {
  settings: defaultSettings,
  isLogin: false,
  userInfo: {
    username: '',
    nickname: '',
    avatar_url: '',
    balance: 0,

    permissions: {},
  },
  userLoading: false,
  baxios: null,
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
      const { userInfo = state.userInfo, userLoading = state.userLoading, isLogin = state.isLogin } = action.payload;
      return {
        ...state,
        userInfo,
        userLoading,
        isLogin,
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
