export type ChatRoom = {
    chat_room_id: string;
    chat_room_name: string;
    chat_room_type: string;
};

export type ChatRoomMember = {
    chat_room_id: string;
    user_id: string;
};

export type Message = {
    message_id: string;
    chat_room_id: string;
    user_id: string;
    content: string;
    timestamp: number;
};

export type OpenConversations = {
    user_id: string;
    chat_room_id: string;
};
