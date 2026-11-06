import React from 'react';

const FileUpload = ({ onUpload, showOptions = true }) => {
  return (
    <div className="file-upload">
      <input type="file" onChange={onUpload} />
      {showOptions && (
        <div className="options">
          <label>Language: <select><option>JavaScript</option></select></label>
          <label>Lines: <input type="checkbox" /></label>
        </div>
      )}
    </div>
  );
};

export default FileUpload;