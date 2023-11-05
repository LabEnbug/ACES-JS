import baxios from "@/utils/getaxios";

export default async function FetchUserInfo(dispatch=null) {
  if (dispatch === null) return;

  dispatch({
    type: 'update-userInfo',
    payload: { userLoading: true },
  });

  try {
    const response = await baxios
      .post('/v1-api/v1/user/info');
    const data = response.data;
    if (data.status !== 200) {
      console.error(data.err_msg);
      dispatch({
        type: 'update-userInfo',
        payload: {userLoading: false, isLogin: false},
      });
      return;
    }
    let userInfo = data.data.user;
    dispatch({
      type: 'update-userInfo',
      payload: {userInfo: userInfo, userLoading: false, isLogin: true},
    });
  } catch (error) {
    console.error(error);
    dispatch({
      type: 'update-userInfo',
      payload: {userLoading: false, isLogin: false},
    });
  }
}

export async function UpdateUserInfoOnly(dispatch=null) {
  if (dispatch === null) return;

  try {
    const response = await baxios
      .post('/v1-api/v1/user/info');
    const data = response.data;
    if (data.status !== 200) {
      console.error(data.err_msg);
      return;
    }
    let userInfo = data.data.user;
    dispatch({
      type: 'update-userInfo',
      payload: {userInfo: userInfo},
    });
  } catch (error) {
    console.error(error);
  }
}