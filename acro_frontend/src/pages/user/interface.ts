
export interface UserCard {
  user_id?: number;
  username?: string;
  nickname?: string;
  reg_time?: string;
}

export interface VideoCard {
  content?: string;
  keyword?: string;
  like_count?: number;
  type?: number;
  user?: UserCard;
  upload_time?: string;
  is_user_like: boolean;
  is_user_favorite: boolean;
  is_user_uploaded: boolean;
  is_user_history: boolean;
  is_user_last_play: boolean;
  cover_url?: string;
  play_url?: string;
  video_uid?: string;
}