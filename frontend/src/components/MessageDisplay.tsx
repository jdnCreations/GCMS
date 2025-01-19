interface MessageProps {
  updateMsg?: string | null;
  errorMsg?: string | null;
}

const MessageDisplay = ({ errorMsg, updateMsg }: MessageProps) => {
  const message = errorMsg || updateMsg || '';
  return (
    <div
      className={`p-1 px-2 rounded my-1 ${errorMsg && 'bg-nook-rose'} ${
        updateMsg && 'bg-nook-light-olive'
      }`}
    >
      {message}
    </div>
  );
};

export default MessageDisplay;
