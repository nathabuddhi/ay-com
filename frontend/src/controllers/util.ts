export function processContent(content: string): string {
    let processed = content.replace(
        /@(\w+)/g,
        "<a href='/profile/$1' class='mention'>@$1</a>"
    );

    processed = processed.replace(
        /#(\w+)/g,
        "<a href='/explore?q=%23$1' class='hashtag'>#$1</a>"
    );

    return processed;
}
