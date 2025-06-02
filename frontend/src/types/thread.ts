export interface Media {
    id: string;
    postId: string;
    url: string;
    type: "image" | "gif" | "video";
    timestamp: string;
}

export interface ThreadResponse {
    threads: Thread[];
}

export interface ThreadDetailResponse {
    thread: Thread;
    replies: ThreadReply[];
}

export interface Thread {
    thread_id: string;
    user_id: string;
    content: string;
    category: string;
    is_advertisement: boolean;
    reply_permission: string;
    pinned: boolean;
    media: Media[];
    poll_options: PollOption[];
    like_count: number;
    repost_count: number;
    reply_count: number;
    posted_at: string;
    is_private: boolean;
    is_liking: boolean;
    is_bookmarking: boolean;
    is_reposting: boolean;
}

export interface Media {
    id: string;
    media_url: string;
    media_type: string;
}

export interface PollOption {
    option: string;
    vote_count: number;
    is_voting: boolean;
}

export interface ThreadLike {
    id: string;
    user_id: string;
}

export interface ThreadRepost {
    id: string;
    user_id: string;
    text: string;
}

export interface ThreadReply {
    id: string;
    user_id: string;
    content: string;
    is_pinned: boolean;
    timestamp: string;
}
