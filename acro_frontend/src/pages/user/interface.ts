
export interface UserCard {
  user_id?: number;
  username?: string;
  nickname?: string;
  reg_time?: string;

  follow_count?: number;
  be_followed?: boolean;
  be_followed_count?: number;

  be_liked_count?: number;
  be_favorite_count?: number;
  be_commented_count?: number;
  be_forwarded_count?: number;
  be_watched_count?: number;
}

export interface VideoCard {
  user?: UserCard;
  video_uid?: string;
  type?: number;
  content?: string;
  keyword?: string;
  upload_time?: string;

  be_liked_count?: number;
  be_favorite_count?: number;
  be_commented_count?: number;
  be_forwarded_count?: number;
  be_watched_count?: number;

  is_user_liked: boolean;
  is_user_favorite: boolean;
  is_user_uploaded: boolean;
  is_user_watched: boolean;
  is_user_last_play: boolean;
  cover_url?: string;
  play_url?: string;
}