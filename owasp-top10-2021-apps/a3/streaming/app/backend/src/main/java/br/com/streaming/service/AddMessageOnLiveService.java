package br.com.streaming.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import br.com.streaming.dto.MessageDTO;
import br.com.streaming.model.Live;
import br.com.streaming.model.Message;
import br.com.streaming.model.User;
import br.com.streaming.repository.LiveRepository;
import br.com.streaming.repository.MessageRepository;
import br.com.streaming.repository.UserRepository;

@Service
public class AddMessageOnLiveService {
	
	@Autowired
	private LiveRepository liveRepository;

	@Autowired
	private UserRepository userRepository;
	
	@Autowired
	private MessageRepository messageRepository;

	public Message add(Long id, MessageDTO messageDto) {
		Live live = liveRepository.findById(id).orElse(null);
		User user = userRepository.findByUsername(messageDto.username).orElse(null);
		
		if(user == null) {
			// Prevent client-controlled username creation/impersonation.
			// Use a generic anonymous account for unauthenticated or unknown users.
			final String ANONYMOUS_USERNAME = "anonymous";
			user = userRepository.findByUsername(ANONYMOUS_USERNAME).orElse(null);
			if (user == null) {
				user = new User(ANONYMOUS_USERNAME, "Anonymous");
				userRepository.save(user);
			}
		}

		Message message = new Message(user, messageDto.content);
		messageRepository.save(message);

		live.addMessage(message);
		liveRepository.save(live);
		return message;
	}
}
